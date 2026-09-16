# Repositori Praktikum Pemrograman Backend

Repositori ini berisi tugas Praktikum Pemrograman Backend Lanjut.

## Struktur Repositori

| Folder | Deskripsi | Pertemuan |
|--------|-----------|-----------|
| `latihan-syntax/` | Latihan sintaks dasar Go | TM1 |
| `latihan-dasar/` | Tugas mandiri: variabel, pointer, struct | TM1 |
| `latihan-fiber/` | Latihan web framework Fiber + database + auth | TM1 → TM5 |
| `api-students/` | REST API mahasiswa dengan autentikasi & JWT | TM2 → TM5 |

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
cp .env.example .env    # isi kredensial database dan JWT
go run .
```

### API Students

#### 1. Siapkan PostgreSQL

Pastikan PostgreSQL sudah terpasang dan berjalan, lalu buat database:

```bash
psql -U postgres -c "CREATE DATABASE praktikum_backend;"
```

#### 2. Jalankan migrasi

Jalankan file SQL secara berurutan untuk membuat tabel mahasiswa, pengguna, dan refresh token:

```bash
psql -U postgres -d praktikum_backend -f api-students/migrations/001_create_students.sql
psql -U postgres -d praktikum_backend -f api-students/migrations/002_create_users.sql
psql -U postgres -d praktikum_backend -f api-students/migrations/003_auth.sql
```

#### 3. Konfigurasi environment

```bash
cd api-students
cp .env.example .env
```

Isi file `.env` dengan kredensial database dan rahasia JWT Anda (buat secret dengan `openssl rand -hex 32`):

```env
APP_PORT=3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=praktikum_backend
DB_SSLMODE=disable
DB_MAX_CONNS=10

JWT_SECRET=your_32_bytes_hex_secret
JWT_ISSUER=api-students
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_DAYS=7
ALLOWED_ORIGINS=http://localhost:5173
```

#### 4. Jalankan aplikasi

```bash
cd api-students
go run .
```

Akses `http://localhost:3000/api/v1/health` untuk memastikan server dan database terhubung.

#### 5. Jalankan unit test

```bash
cd api-students
go test ./app/service/... -v
```

---

## Struktur Proyek (Clean Architecture — TM4 & TM5)

Mulai pertemuan 4 dan 5, proyek `api-students` dan `latihan-fiber` direstrukturisasi mengikuti prinsip Clean Architecture dengan penambahan modul autentikasi dan keamanan:

```
api-students/
├── app/
│   ├── model/          struct entitas (student, user, auth) dan request/response
│   ├── repository/     akses database tabel students, users, dan refresh_tokens
│   └── service/        business rules murni (student & auth rules) serta service
├── config/             app.go, env.go, logger.go
├── database/           koneksi pool PostgreSQL (pgxpool)
├── helper/             security (bcrypt), jwt, context, request parser, response envelope
├── logs/               output rotasi log (tidak di-commit)
├── middleware/         middleware global, RequireJSON, RequireAuth, LoginRateLimit
├── migrations/         file skrip migrasi DDL SQL (001, 002, 003)
├── route/              pendaftaran route (publik, auth, students terproteksi)
├── .env                (tidak di-commit)
├── .env.example        template environment variable
└── main.go             inisialisasi dependensi, JWT manager, dan graceful shutdown
```

### Pemetaan ke Layer Clean Architecture

| Layer | Folder | Keterangan |
|-------|--------|------------|
| 1. Entities | `app/model/` | `Student`, `User`, `RefreshToken`, struct request & response |
| 2. Use Cases | `app/service/` | Business rules murni (`student_rules.go`, `auth_rules.go`) |
| 3. Interface Adapters | `app/repository/`, `helper/`, `app/service/` | Repository (gateway), presenter (response helper), service (controller) |
| 4. Frameworks & Drivers | `config/`, `database/`, `middleware/`, `route/`, `main.go` | Fiber web framework, PostgreSQL driver, routing, main server |

---

## Skema Database

### 1. Tabel Students
```sql
CREATE TABLE IF NOT EXISTS students (
    id         SERIAL       PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    nim        VARCHAR(50)  NOT NULL,
    grade      FLOAT        NOT NULL DEFAULT 0,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS students_nim_lower_key ON students (LOWER(nim));
CREATE INDEX IF NOT EXISTS students_name_lower_idx ON students (LOWER(name));
```

### 2. Tabel Users (TM5)
```sql
CREATE TABLE IF NOT EXISTS users (
    id         SERIAL       PRIMARY KEY,
    username   VARCHAR(50)  NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(20)  NOT NULL DEFAULT 'user',
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_username_lower_key ON users (LOWER(username));
CREATE INDEX IF NOT EXISTS users_email_lower_idx ON users (LOWER(email));
```

### 3. Tabel Refresh Tokens (TM5)
```sql
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    INTEGER     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS refresh_tokens_user_id_idx ON refresh_tokens (user_id);
```

---

## Environment Variables

| Variabel | Keterangan | Contoh |
|----------|------------|--------|
| `APP_NAME` | Nama aplikasi | `Praktikum Backend Lanjut` |
| `APP_PORT` | Port aplikasi | `3000` |
| `LOG_LEVEL` | Level log: `debug`, `info`, `warn`, `error` | `info` |
| `DB_HOST` | Host PostgreSQL | `localhost` |
| `DB_PORT` | Port PostgreSQL | `5432` |
| `DB_USER` | Username database | `postgres` |
| `DB_PASSWORD` | Password database | `secret` |
| `DB_NAME` | Nama database | `praktikum_backend` |
| `DB_SSLMODE` | Mode SSL | `disable` |
| `DB_MAX_CONNS` | Maksimum koneksi pool | `10` |
| `JWT_SECRET` | Kunci rahasia HMAC-SHA256 untuk signing JWT | *(32-byte hex random string)* |
| `JWT_ISSUER` | Nama penerbit token JWT | `api-students` |
| `JWT_ACCESS_TTL_MINUTES` | Masa berlaku Access Token (menit) | `15` |
| `JWT_REFRESH_TTL_DAYS` | Masa berlaku Refresh Token (hari) | `7` |
| `ALLOWED_ORIGINS` | Asal domain CORS yang diizinkan | `http://localhost:5173` |

> File `.env` **tidak boleh di-commit**. Gunakan `.env.example` sebagai referensi.

---

## Kontrak API

### Endpoint Publik & Health Check

| Metode | Endpoint | Keterangan | Auth | Status |
|--------|----------|------------|------|--------|
| `GET` | `/api/v1/health` | Pengecekan koneksi server & database | Publik | `200`, `503` |

### Endpoint Autentikasi (`/api/v1/auth`)

| Metode | Endpoint | Body Permintaan | Keterangan | Status yang Mungkin |
|--------|----------|-----------------|------------|---------------------|
| `POST` | `/api/v1/auth/register` | `{"username": "...", "email": "...", "password": "..."}` | Pendaftaran akun baru (role otomatis `'user'`) | `201`, `400`, `409`, `415`, `422` |
| `POST` | `/api/v1/auth/login` | `{"username": "...", "password": "..."}` | Login untuk memperoleh Access & Refresh Token (max 5 percobaan/menit) | `200`, `400`, `401`, `403`, `415`, `422`, `429` |
| `POST` | `/api/v1/auth/refresh` | `{"refresh_token": "..."}` | Rotasi token (menerbitkan pair baru dan mencabut token lama) | `200`, `400`, `401`, `415` |
| `POST` | `/api/v1/auth/logout` | `{"refresh_token": "..."}` | Mencabut refresh token yang aktif | `200`, `400`, `415` |
| `GET` | `/api/v1/auth/me` | *(Kosong)* | Melihat data profil pengguna yang sedang login | `200`, `401` *(Wajib Bearer Token)* |

### Endpoint Mahasiswa (`/api/v1/students` — Wajib Autentikasi)

> **Catatan Keamanan:** Seluruh endpoint `/api/v1/students` wajib menyertakan header HTTP:  
> `Authorization: Bearer <access_token>`

| Metode | Endpoint | Parameter | Contoh Body | Status yang Mungkin |
|--------|----------|-----------|-------------|---------------------|
| `GET` | `/api/v1/students` | **Query**: `page`, `limit`, `search`, `sort`, `order`, `is_active` | *(Kosong)* | `200`, `401`, `500` |
| `GET` | `/api/v1/students/:id` | **Path**: `id` | *(Kosong)* | `200`, `400`, `401`, `404` |
| `POST` | `/api/v1/students` | *(Kosong)* | `{"nim": "0431", "name": "Budi", "grade": 90, "is_active": true}` | `201`, `400`, `401`, `409`, `415`, `422` |
| `PUT` | `/api/v1/students/:id` | **Path**: `id` | `{"nim": "0431", "name": "Budi", "grade": 95, "is_active": true}` | `200`, `400`, `401`, `404`, `409`, `415`, `422` |
| `PATCH` | `/api/v1/students/:id` | **Path**: `id` | `{"grade": 95}` | `200`, `400`, `401`, `404`, `409`, `415`, `422` |
| `DELETE` | `/api/v1/students/:id` | **Path**: `id` | *(Kosong)* | `204`, `400`, `401`, `404` |

---

*Lihat [AI-USAGE.md](./AI-USAGE.md) untuk detail bantuan penggunaan AI dan referensi eksternal.*
