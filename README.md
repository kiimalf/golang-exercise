# Repositori Praktikum Pemrograman Backend

Repositori ini berisi kumpulan tugas Praktikum Pemrograman Backend Lanjut, disusun per Modul (Tatap Muka / TM).

---

## [TM 1] Persiapan Lingkungan & Sintaks Go dengan Fiber

**Latihan Panduan:**

- Latihan web Fiber:
  ```bash
  cd TM1/latihan-fiber
  go run main.go
  ```
- Latihan sintaks dasar:
  ```bash
  cd TM1/latihan-syntax
  go run main.go
  ```

**Tugas Mandiri:**

- Tugas 2 (Variabel):
  ```bash
  cd TM1/tugas_mandiri/tugas2_variabel
  go run main.go
  ```
- Tugas 3 (Pointer):
  ```bash
  cd TM1/tugas_mandiri/tugas3_pointer
  go run main.go
  ```
- Tugas 4 (Struct):
  ```bash
  cd TM1/tugas_mandiri/tugas4_struct
  go run main.go
  ```

---

## [TM 2] REST API & HTTP Deep Dive

**Tugas Mandiri:**

- API Students:
  ```bash
  cd TM2/api-students
  go run .
  ```

- Dokumen Kontrak API:

  | Metode | Endpoint | Parameter | Contoh Body Permintaan | Status yang Mungkin | Contoh Respons |
  |---|---|---|---|---|---|
  | `GET` | `/api/v1/students` | **Query**: `page`, `limit`, `search`, `sort`, `order`, `is_active` | *(Kosong)* | `200` | `{"success": true, "message": "daftar mahasiswa berhasil diambil", "data": [...], "meta": {"page": 1, "limit": 10, "total": 5, "total_pages": 1}}` |
  | `GET` | `/api/v1/students/:id` | **Path**: `id` | *(Kosong)* | `200`, `400`, `404` | `{"success": true, "message": "mahasiswa ditemukan", "data": {"id": 1, "nim": "0431", "name": "Budi", "grade": 90, "is_active": true}}` |
  | `POST` | `/api/v1/students` | *(Kosong)* | `{"nim": "0431", "name": "Budi", "grade": 90, "is_active": true}` | `201`, `400`, `409`, `415`, `422` | `{"success": true, "message": "mahasiswa berhasil dibuat", "data": {"id": 1, "nim": "0431", "name": "Budi", "grade": 90, "is_active": true}}` |
  | `PUT` | `/api/v1/students/:id` | **Path**: `id` | `{"nim": "0431", "name": "Budi", "grade": 95, "is_active": true}` | `200`, `400`, `404`, `409`, `415`, `422` | `{"success": true, "message": "mahasiswa berhasil diganti seluruhnya", "data": {"id": 1, "nim": "0431", "name": "Budi", "grade": 95, "is_active": true}}` |
  | `PATCH` | `/api/v1/students/:id` | **Path**: `id` | `{"grade": 95}` | `200`, `400`, `404`, `409`, `415`, `422` | `{"success": true, "message": "mahasiswa berhasil diperbarui sebagian", "data": {"id": 1, "nim": "0431", "name": "Budi", "grade": 95, "is_active": true}}` |
  | `DELETE` | `/api/v1/students/:id` | **Path**: `id` | *(Kosong)* | `204`, `400`, `404` | *(Tidak ada body respons)* |

## [TM 3] *(Akan Datang)*

---

*Lihat [AI-USAGE.md](./AI-USAGE.md) untuk detail bantuan penggunaan AI dan referensi eksternal.*
