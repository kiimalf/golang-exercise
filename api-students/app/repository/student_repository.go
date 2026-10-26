package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"api-students/app/model"
)

var (
	ErrNotFound  = errors.New("Data tidak ditemukan")
	ErrDuplicate = errors.New("Data Sudah ada")
)

type StudentRepository interface {
	WithTx(tx DBTX) StudentRepository
	LockStudent(ctx context.Context, studentID int) error

	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindAfterCursor(ctx context.Context, q model.CursorQuery) ([]model.Student, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, s model.Student) (model.Student, error)
	Update(ctx context.Context, s model.Student) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

var kolomUrut = map[string]string{
	"id":    "id",
	"name":  "name",
	"nim":   "nim",
	"grade": "grade",
}

const studentColumns = `id, owner_id, name, nim, grade, is_active, created_at`

type studentPostgresRepository struct {
	db DBTX
}

func NewStudentRepository(db DBTX) StudentRepository {
	return &studentPostgresRepository{db: db}
}

func scanStudent(row pgx.Row) (model.Student, error) {
	var s model.Student
	err := row.Scan(&s.ID, &s.OwnerID, &s.Name, &s.NIM, &s.Grade, &s.IsActive, &s.CreatedAt)
	return s, err
}

func buildFilter(q model.ListQuery) (string, []any) {
	where := " WHERE 1 = 1"
	args := []any{}

	if q.Search != "" {
		where += fmt.Sprintf(" AND name ILIKE $%d",
			len(args)+1)
		args = append(args, "%"+q.Search+"%")
	}

	if q.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", len(args)+1)
		args = append(args, *q.IsActive)
	}

	return where, args
}

func (r *studentPostgresRepository) WithTx(tx DBTX) StudentRepository {
	return &studentPostgresRepository{db: tx}
}

func (r *studentPostgresRepository) LockStudent(ctx context.Context, studentID int) error {
	var id int
	err := r.db.QueryRow(ctx,
		`SELECT id FROM students 
		WHERE id = $1 FOR UPDATE`, studentID,
	).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("Mengunci data mahasiswa: %w", err)
	}
	return nil
}

func (r *studentPostgresRepository) FindAll(
	ctx context.Context, q model.ListQuery,
) ([]model.Student, int, error) {
	where, args := buildFilter(q)

	var total int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM students"+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("Menghitung student: %w", err)
	}

	arah := "ASC"
	if q.Order == "desc" {
		arah = "DESC"
	}

	sqlText := fmt.Sprintf(
		`SELECT `+studentColumns+`
		FROM students%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		where, kolomUrut[q.Sort], arah, len(args)+1, len(args)+2,
	)
	args = append(args, q.Limit, q.Offset())

	rows, err := r.db.Query(ctx, sqlText, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Mengambil daftar student: %w", err)
	}
	defer rows.Close()

	hasil := []model.Student{}
	for rows.Next() {
		s, err := scanStudent(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("Membaca baris student: %w", err)
		}
		hasil = append(hasil, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("Membaca hasil query: %w", err)
	}

	return hasil, total, nil
}

func (r *studentPostgresRepository) FindAfterCursor(
	ctx context.Context, q model.CursorQuery,
) ([]model.Student, error) {
	args := []any{}
	where := " WHERE 1 = 1"

	if q.Search != "" {
		args = append(args, "%"+q.Search+"%")
		where += fmt.Sprintf(" AND name ILIKE $%d", len(args))
	}
	if q.IsActive != nil {
		args = append(args, *q.IsActive)
		where += fmt.Sprintf(" AND is_active = $%d", len(args))
	}
	if q.After != nil {
		args = append(args, q.After.CreatedAt, q.After.ID)
		where += fmt.Sprintf(" AND (created_at, id) < ($%d, $%d)",
			len(args)-1, len(args))
	}

	args = append(args, q.Limit+1)
	query := fmt.Sprintf(
		`SELECT `+studentColumns+`
		FROM students%s
		ORDER BY created_at DESC, id DESC
		LIMIT $%d`,
		where, len(args),
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Mengambil data student: %w", err)
	}
	defer rows.Close()

	result := []model.Student{}
	for rows.Next() {
		s, err := scanStudent(rows)
		if err != nil {
			return nil, fmt.Errorf("Membaca row student: %w", err)
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Membaca hasil query student: %w", err)
	}

	return result, nil
}

func (r *studentPostgresRepository) FindByID(
	ctx context.Context, id int,
) (model.Student, error) {
	s, err := scanStudent(r.db.QueryRow(ctx,
		"SELECT "+studentColumns+" FROM students WHERE id = $1", id,
	))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("Mengambil student: %w", err)
	}

	return s, nil
}

func (r *studentPostgresRepository) Create(
	ctx context.Context, s model.Student,
) (model.Student, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO students (owner_id, name, nim, grade, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		s.OwnerID, s.Name, s.NIM, s.Grade, s.IsActive,
	).Scan(&s.ID, &s.CreatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("Menyimpan student: %w", err)
	}

	return s, nil
}

func (r *studentPostgresRepository) Update(
	ctx context.Context, s model.Student,
) (model.Student, error) {
	s, err := scanStudent(r.db.QueryRow(ctx,
		`UPDATE students SET name = $1, nim = $2, grade = $3, is_active = $4
		WHERE id = $5
		RETURNING `+studentColumns,
		s.Name, s.NIM, s.Grade, s.IsActive, s.ID,
	))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("Memperbarui student: %w", err)
	}

	return s, nil
}

func (r *studentPostgresRepository) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM students WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Menghapus student: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
