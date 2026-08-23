// Bagian 1 : Penyimpanan dan bantuan
package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var users []User
var nextID = 1

func findUserIndex(id int) int {
	for i := range users {
		if users[i].ID == id {
			return i
		}
	}
	return -1
}

func cocokPencarian(u User, kata string) bool {
	kata = strings.ToLower(kata)
	return strings.Contains(strings.ToLower(u.Email), kata)
}

func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

// Bagian 2: Daftar dengan saring, urut, dan penggal
func listUsers(c *fiber.Ctx) error {
	q := parseListQuery(c)

	hasil := []User{}
	for _, u := range users {
		if q.IsActive != nil && u.IsActive != *q.IsActive {
			continue
		}
		if q.Search != "" && !cocokPencarian(u, q.Search) {
			continue
		}
		hasil = append(hasil, u)
	}

	sort.SliceStable(hasil, func(i, j int) bool {
		var lebihKecil bool
		switch q.Sort {
		case "username":
			lebihKecil = hasil[i].Username < hasil[j].Username
		case "email":
			lebihKecil = hasil[i].Email < hasil[j].Email
		case "created_at":
			lebihKecil = hasil[i].CreatedAt.Before(hasil[j].CreatedAt)
		default:
			lebihKecil = hasil[i].ID < hasil[j].ID
		}
		if q.Order == "desc" {
			return !lebihKecil
		}
		return lebihKecil
	})

	total := len(hasil)
	totalPages := (total + q.Limit - 1) / q.Limit
	mulai := (q.Page - 1) * q.Limit
	if mulai > total {
		mulai = total
	}
	akhir := mulai + q.Limit
	if akhir > total {
		akhir = total
	}

	return okList(c, "Daftar user berhasil diambil", hasil[mulai:akhir], &Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// Bagian 3: Ambil satu dan tambah baru
func getUser(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	i := findUserIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return ok(c, "User ditemukan", users[i])
}

func createUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" {
		errs["username"] = "Wajib diisi"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "Format email tidak valid"
	}
	if len(req.Password) < 8 {
		errs["password"] = "Minimal 8 karakter"
	}
	for _, u := range users {
		if strings.EqualFold(u.Username, req.Username) {
			errs["username"] = "Sudah dipakai"
		}
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	baru := User{
		ID:        nextID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	users = append(users, baru)
	nextID++

	return created(c, "User berhasil dibuat", baru,
		"/api/v1/users"+strconv.Itoa(baru.ID))
}

// Bagian 4
func replaceUser(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	i := findUserIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "User tidak ditemukan")
	}

	var req ReplaceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "Wajib diisi pada PUT"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "format email tidak valid"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	users[i].Username = req.Username
	users[i].Email = req.Email
	users[i].IsActive = req.IsActive

	return ok(c, "User berhasil diganti seluruhnya", users[i])
}

func patchUser(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	i := findUserIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "User tidak ditemukan")
	}

	var req PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Body harus berupa JSON yang valid")
	}

	if req.Username == nil && req.Email == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "Tidak ada field yang diubah")
	}

	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			return failValidation(c, map[string]string{"Username": "Tidak boleh kosong"})
		}
		users[i].Username = *req.Username
	}
	if req.Email != nil {
		if !strings.Contains(*req.Email, "@") {
			return failValidation(c, map[string]string{"email": "Format email tidak valid"})
		}
		users[i].Email = *req.Email
	}
	if req.IsActive != nil {
		users[i].IsActive = *req.IsActive
	}

	return ok(c, "User berhasil diperbarui sebagian", users[i])
}

func deleteUser(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "Id harus berupa angka positif")
	}

	i := findUserIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "User tidak ditemukan")
	}

	users = append(users[:i], users[i+1:]...)

	return noContent(c)
}
