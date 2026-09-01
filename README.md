# Repositori Praktikum Pemrograman Backend

Repositori ini berisi tugas Praktikum Pemrograman Backend Lanjut.

## Struktur Repositori

| Folder | Deskripsi | Pertemuan |
|--------|-----------|-----------|
| `latihan-syntax/` | Latihan sintaks dasar Go | TM1 |
| `latihan-dasar/` | Tugas mandiri: variabel, pointer, struct | TM1 |
| `latihan-fiber/` | Latihan web framework Fiber + database | TM1 → TM3 |
| `api-students/` | REST API mahasiswa (tugas mandiri) | TM2 → TM3 |

---

## Cara Menjalankan

### Latihan Syntax

```bash
cd latihan-syntax
go run main.go
```

### Latihan Dasar

```bash
cd latihan-dasar/variable && go run main.go
cd latihan-dasar/pointer && go run main.go
cd latihan-dasar/struct && go run main.go
```

### Latihan Fiber

```bash
cd latihan-fiber
cp .env.example .env    # isi kredensial database
go run .
```

### API Students

#### 1. Siapkan PostgreSQL

Pastikan PostgreSQL sudah terpasang dan berjalan, lalu buat database:

```bash
psql -U postgres -c "CREATE DATABASE praktikum_backend;"
```

#### 2. Jalankan migrasi

Jalankan file SQL untuk membuat tabel:

```bash
psql -U postgres -d praktikum_backend -f api-students/migrations/001_create_students.sql
```

#### 3. Konfigurasi environment

```bash
cd api-students
cp .env.example .env
```

Isi file `.env` dengan kredensial database Anda:

```env
APP_PORT=3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=praktikum_backend
DB_SSLMODE=disable
DB_MAX_CONNS=10
```

#### 4. Jalankan aplikasi

```bash
cd api-students
go run .
```

Akses `http://localhost:3000/api/v1/health` untuk memastikan server dan database terhubung.

---

## Skema Tabel Students

```sql
CREATE TABLE IF NOT EXISTS students (
    id         SERIAL       PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    nim        VARCHAR(50)  NOT NULL,
    grade      FLOAT        NOT NULL DEFAULT 0,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- NIM unik tanpa membedakan huruf besar-kecil
CREATE UNIQUE INDEX IF NOT EXISTS students_nim_lower_key
    ON students (LOWER(nim));

-- Indeks untuk pencarian nama
CREATE INDEX IF NOT EXISTS students_name_lower_idx
    ON students (LOWER(name));
```

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | `SERIAL` | Primary key, auto-increment |
| `name` | `VARCHAR(255)` | Nama mahasiswa, wajib diisi |
| `nim` | `VARCHAR(50)` | NIM mahasiswa, wajib unik (case-insensitive) |
| `grade` | `FLOAT` | Nilai mahasiswa, default `0` |
| `is_active` | `BOOLEAN` | Status aktif, default `TRUE` |
| `created_at` | `TIMESTAMPTZ` | Waktu pembuatan, default `NOW()` |

---

## Environment Variables

| Variabel | Keterangan | Contoh |
|----------|------------|--------|
| `APP_PORT` | Port aplikasi | `3000` |
| `DB_HOST` | Host PostgreSQL | `localhost` |
| `DB_PORT` | Port PostgreSQL | `5432` |
| `DB_USER` | Username database | `postgres` |
| `DB_PASSWORD` | Password database | `secret` |
| `DB_NAME` | Nama database | `praktikum_backend` |
| `DB_SSLMODE` | Mode SSL | `disable` |
| `DB_MAX_CONNS` | Maksimum koneksi pool | `10` |

> File `.env` **tidak di-commit**. Gunakan `.env.example` sebagai referensi.

---

## Kontrak API Students

| Metode | Endpoint | Parameter | Contoh Body Permintaan | Status yang Mungkin | Contoh Respons |
|--------|----------|-----------|------------------------|---------------------|----------------|
| `GET` | `/api/v1/students` | **Query**: `page`, `limit`, `search`, `sort`, `order`, `is_active` | *(Kosong)* | `200` | `{"success": true, "message": "daftar mahasiswa berhasil diambil", "data": [...], "meta": {"page": 1, "limit": 10, "total": 5, "total_pages": 1}}` |
| `GET` | `/api/v1/students/:id` | **Path**: `id` | *(Kosong)* | `200`, `400`, `404` | `{"success": true, "message": "mahasiswa ditemukan", "data": {"id": 1, "nim": "0431", "name": "Budi", "grade": 90, "is_active": true}}` |
| `POST` | `/api/v1/students` | *(Kosong)* | `{"nim": "0431", "name": "Budi", "grade": 90, "is_active": true}` | `201`, `400`, `409`, `415`, `422` | `{"success": true, "message": "mahasiswa berhasil dibuat", "data": {"id": 1, "nim": "0431", "name": "Budi", "grade": 90, "is_active": true}}` |
| `PUT` | `/api/v1/students/:id` | **Path**: `id` | `{"nim": "0431", "name": "Budi", "grade": 95, "is_active": true}` | `200`, `400`, `404`, `409`, `415`, `422` | `{"success": true, "message": "mahasiswa berhasil diganti seluruhnya", "data": {"id": 1, "nim": "0431", "name": "Budi", "grade": 95, "is_active": true}}` |
| `PATCH` | `/api/v1/students/:id` | **Path**: `id` | `{"grade": 95}` | `200`, `400`, `404`, `409`, `415`, `422` | `{"success": true, "message": "mahasiswa berhasil diperbarui sebagian", "data": {"id": 1, "nim": "0431", "name": "Budi", "grade": 95, "is_active": true}}` |
| `DELETE` | `/api/v1/students/:id` | **Path**: `id` | *(Kosong)* | `204`, `400`, `404` | *(Tidak ada body respons)* |
| `GET` | `/api/v1/health` | *(Kosong)* | *(Kosong)* | `200`, `503` | `{"success": true, "message": "Server berjalan"}` |

---

*Lihat [AI-USAGE.md](./AI-USAGE.md) untuk detail bantuan penggunaan AI dan referensi eksternal.*
