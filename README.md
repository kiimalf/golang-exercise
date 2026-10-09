# Repositori Praktikum Pemrograman Backend

Repositori ini berisi tugas Praktikum Pemrograman Backend Lanjut.

## Struktur Repositori

| Folder | Deskripsi | Pertemuan |
|--------|-----------|-----------|
| `latihan-syntax/` | Latihan sintaks dasar Go | TM1 |
| `latihan-dasar/` | Tugas mandiri: variabel, pointer, struct | TM1 |
| `latihan-fiber/` | Latihan web framework Fiber + database + autentikasi & otorisasi RBAC | TM1 → TM6 |
| `api-students/` | REST API mahasiswa dengan autentikasi JWT & otorisasi RBAC (ownership) | TM2 → TM6 |

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

Jalankan file SQL secara berurutan untuk membuat tabel mahasiswa, pengguna, refresh token, skema RBAC, dan kebijakan hak akses mahasiswa:

```bash
psql -U postgres -d praktikum_backend -f api-students/migrations/001_create_students.sql
psql -U postgres -d praktikum_backend -f api-students/migrations/002_create_users.sql
psql -U postgres -d praktikum_backend -f api-students/migrations/003_auth.sql
psql -U postgres -d praktikum_backend -f api-students/migrations/004_rbac.sql
psql -U postgres -d praktikum_backend -f api-students/migrations/005_student_permissions.sql
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

## Struktur Proyek (Clean Architecture — TM4 s.d. TM6)

Proyek `api-students` dan `latihan-fiber` mengadopsi prinsip Clean Architecture yang dilengkapi dengan layer autentikasi JWT dan otorisasi Role-Based Access Control (RBAC):

```
api-students/
├── app/
│   ├── model/          struct entitas (student, user, auth) dan request/response
│   ├── repository/     akses database tabel students, users, roles, dan refresh_tokens
│   └── service/        business rules murni (student, auth, authz rules) serta service
├── config/             app.go, env.go, logger.go
├── database/           koneksi pool PostgreSQL (pgxpool)
├── helper/             security (bcrypt), jwt, authz (PermissionSet), context, request, response envelope
├── logs/               output rotasi log (tidak di-commit)
├── middleware/         middleware global, RequireJSON, RequireAuth, RequirePermission, LoginRateLimit
├── migrations/         file skrip migrasi DDL SQL (001 s.d. 005)
├── route/              pendaftaran route (publik, auth, users, students terproteksi)
├── .env                (tidak di-commit)
├── .env.example        template environment variable
└── main.go             inisialisasi dependensi, PermissionSet cache, JWT manager, dan graceful shutdown
```

### Pemetaan ke Layer Clean Architecture

| Layer | Folder | Keterangan |
|-------|--------|------------|
| 1. Entities | `app/model/` | `Student`, `User`, `RefreshToken`, struct request & response |
| 2. Use Cases | `app/service/` | Business & authorization rules murni (`student_rules.go`, `auth_rules.go`, `authz_rules.go`, `student_authz_rules.go`) |
| 3. Interface Adapters | `app/repository/`, `helper/`, `app/service/` | Repository (gateway database), presenter (response helper), authorization helper (`PermissionSet`), service (controller) |
| 4. Frameworks & Drivers | `config/`, `database/`, `middleware/`, `route/`, `main.go` | Fiber web framework, PostgreSQL driver (pgxpool), routing & middleware penjaga akses, main server |

---

## Skema Database

### 1. Tabel Students (TM2, TM5, & TM6)
```sql
CREATE TABLE IF NOT EXISTS students (
    id         SERIAL       PRIMARY KEY,
    owner_id   INTEGER      REFERENCES users(id),
    name       VARCHAR(255) NOT NULL,
    nim        VARCHAR(50)  NOT NULL,
    grade      FLOAT        NOT NULL DEFAULT 0,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS students_nim_lower_key ON students (LOWER(nim));
CREATE INDEX IF NOT EXISTS students_name_lower_idx ON students (LOWER(name));
CREATE INDEX IF NOT EXISTS students_owner_id_idx ON students (owner_id);
```

### 2. Tabel Users (TM5 & TM6)
```sql
CREATE TABLE IF NOT EXISTS users (
    id         SERIAL       PRIMARY KEY,
    username   VARCHAR(50)  NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(20)  NOT NULL DEFAULT 'user' REFERENCES roles(name) ON UPDATE CASCADE,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_username_lower_key ON users (LOWER(username));
CREATE INDEX IF NOT EXISTS users_email_lower_idx ON users (LOWER(email));
CREATE INDEX IF NOT EXISTS users_role_idx ON users (role);
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

### 4. Tabel Roles (TM6 — RBAC)
```sql
CREATE TABLE IF NOT EXISTS roles (
    name        VARCHAR(20)  PRIMARY KEY,
    description VARCHAR(150) NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

### 5. Tabel Permissions (TM6 — RBAC)
```sql
CREATE TABLE IF NOT EXISTS permissions (
    name        VARCHAR(50)  PRIMARY KEY,
    description VARCHAR(150) NOT NULL
);
```

### 6. Tabel Role Permissions (TM6 — Relasi RBAC)
```sql
CREATE TABLE IF NOT EXISTS role_permissions (
    role_name       VARCHAR(20) NOT NULL REFERENCES roles(name) ON DELETE CASCADE,
    permission_name VARCHAR(50) NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, permission_name)
);
```

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
| `GET` | `/api/v1/auth/me` | *(Kosong)* | Melihat data profil pengguna login beserta daftar `permissions` role miliknya | `200`, `401` *(Wajib Bearer Token)* |

### Endpoint Manajemen Pengguna (`/api/v1/users` — TM6 RBAC)

> **Catatan:** Seluruh endpoint `/api/v1/users` wajib menyertakan header `Authorization: Bearer <access_token>`.

| Metode | Endpoint | Parameter / Body | Keterangan & Kebijakan Hak Akses | Status yang Mungkin |
|--------|----------|------------------|----------------------------------|---------------------|
| `GET` | `/api/v1/users` | **Query**: `page`, `limit`, `search`, `sort`, `order`, `is_active` | Melihat daftar seluruh pengguna (Izin: `user:list` — Admin, Staff) | `200`, `401`, `403`, `500` |
| `POST` | `/api/v1/users` | `{"username": "...", "email": "...", "password": "..."}` | Menambahkan pengguna baru (Izin: `user:create` — Admin) | `201`, `400`, `401`, `403`, `409`, `415`, `422` |
| `GET` | `/api/v1/users/:id` | **Path**: `id` | Melihat profil pengguna tertentu (Milik sendiri atau Izin `user:read:any`) | `200`, `400`, `401`, `403`, `404` |
| `PUT` | `/api/v1/users/:id` | `{"username": "...", "email": "...", "is_active": true}` | Mengganti data pengguna secara penuh (Milik sendiri atau Izin `user:update:any`) | `200`, `400`, `401`, `403`, `404`, `409`, `415`, `422` |
| `PATCH` | `/api/v1/users/:id` | `{"username": "...", "email": "..."}` | Mengubah sebagian data pengguna (Milik sendiri atau Izin `user:update:any`) | `200`, `400`, `401`, `403`, `404`, `409`, `415`, `422` |
| `DELETE` | `/api/v1/users/:id` | **Path**: `id` | Menghapus pengguna (Izin: `user:delete` — Admin, dilarang hapus akun sendiri) | `204`, `400`, `401`, `403`, `404` |
| `PATCH` | `/api/v1/users/:id/role` | `{"role": "admin\|staff\|user"}` | Mengubah role pengguna lain (Izin: `role:assign` — Admin, dilarang ubah role sendiri) | `200`, `400`, `401`, `403`, `404`, `415`, `422` |

### Endpoint Mahasiswa (`/api/v1/students` — TM6 RBAC & Ownership)

> **Catatan Keamanan:** Seluruh endpoint `/api/v1/students` wajib menyertakan header HTTP:  
> `Authorization: Bearer <access_token>`

| Metode | Endpoint | Parameter / Body | Keterangan & Kebijakan Hak Akses | Status yang Mungkin |
|--------|----------|------------------|----------------------------------|---------------------|
| `GET` | `/api/v1/students` | **Query**: `page`, `limit`, `search`, `sort`, `order`, `is_active` | Melihat daftar mahasiswa (Izin: `student:list` — Admin, Staff) | `200`, `401`, `403`, `500` |
| `GET` | `/api/v1/students/:id` | **Path**: `id` | Melihat detail mahasiswa (Milik sendiri berdasarkan `owner_id` atau Izin `student:read:any` — Admin, Staff) | `200`, `400`, `401`, `403`, `404` |
| `POST` | `/api/v1/students` | `{"nim": "0431", "name": "Budi", "grade": 90, "is_active": true}` | Mendaftarkan mahasiswa baru; `owner_id` diisi otomatis dari token (Izin: `student:create` — Admin, Staff) | `201`, `400`, `401`, `403`, `409`, `415`, `422` |
| `PUT` | `/api/v1/students/:id` | `{"nim": "0431", "name": "Budi", "grade": 95, "is_active": true}` | Mengubah data mahasiswa penuh (Milik sendiri berdasarkan `owner_id` atau Izin `student:update:any` — Admin) | `200`, `400`, `401`, `403`, `404`, `409`, `415`, `422` |
| `PATCH` | `/api/v1/students/:id` | `{"grade": 95}` | Mengubah data mahasiswa sebagian (Milik sendiri berdasarkan `owner_id` atau Izin `student:update:any` — Admin) | `200`, `400`, `401`, `403`, `404`, `409`, `415`, `422` |
| `DELETE` | `/api/v1/students/:id` | **Path**: `id` | Menghapus mahasiswa (Izin: `student:delete` — Admin) | `204`, `400`, `401`, `403`, `404` |

---

## Mekanisme Keamanan Otorisasi (TM6)

1. **Prinsip Fail Closed**: Nilai izin *default* selalu `false`. Apabila token tidak dikenali, role tidak valid, atau permission tidak terdaftar, akses langsung ditolak (`401` atau `403`).
2. **In-Memory PermissionSet Cache**: Pemetaan role dan permissions dimuat sekali saat aplikasi dijalankan melalui `app/repository/role_repository.go` dan disimpan dalam struktur `map[string]map[string]struct{}` untuk efisiensi pencarian $O(1)$ tanpa membebani database pada setiap request.
3. **Pemisahan Penjaga (Route Guard vs Service Ownership)**:
   - Keputusan non-data diletakkan di middleware (`RequirePermission`).
   - Keputusan berbasis data (`owner_id` entitas) diletakkan di layer service (`CanAccessUser` dan `CanAccessStudent`) untuk mencegah kebocoran informasi dan *timing attack*.
4. **Perlindungan Mass Assignment**: Kolom `owner_id` dan `role` diisolasi dari DTO request biasa dan hanya dapat diisi melalui token JWT atau endpoint khusus terproteksi.
5. **Audit Trail Logging**: Middleware `RequestLogger` mencatat atribut terstruktur `user_id` dan `role` pada setiap request yang lolos autentikasi untuk kemudahan penelusuran akses.

---

*Lihat [AI-USAGE.md](./AI-USAGE.md) untuk detail bantuan penggunaan AI dan referensi eksternal.*
