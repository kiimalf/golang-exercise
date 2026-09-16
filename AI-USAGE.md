# Laporan Penggunaan Bantuan AI dan Referensi Eksternal

Dokumen ini berisi daftar bantuan AI dan referensi eksternal yang digunakan selama pengerjaan Tugas Mandiri Praktikum Pemrograman Backend.

## Pertemuan 1 (TM 1)

### 🤖 Bantuan AI
Dalam pengerjaan TM 1, AI digunakan sebagai tutor interaktif untuk membantu menjelaskan konsep dan mereviu kode, dengan rincian sebagai berikut:
- **Penjelasan Modul:** Membantu merangkum isi modul dan memberikan penjelasan tambahan tentang konsep Go yang krusial bagi pemula tetapi tidak ada di modul (seperti aturan *Exported/Unexported variables* dengan huruf kapital, dan konsep dasar *Error Handling* di Go).
- **Pemahaman Sintaks:** Membantu menjelaskan cara kerja idiom *"Comma Ok"* pada pengecekan Map (`if n, ada := skor["Budi"]; ada`).
- **Pemahaman Konsep Pointer & Slice:** Membantu memberikan analogi dan menjelaskan alur kerja (*flow*) fungsi `swap` menggunakan *pointer*, serta menjelaskan mengapa memodifikasi *slice* (seperti operasi `append`) memerlukan *pointer* di parameter fungsi `updateSlice`.
- **Pemahaman Teori:** Menjelaskan perbedaan fundamental antara konsep *Pass by Value* dan *Pass by Reference*.
- **Review Kode:** AI digunakan untuk memverifikasi fungsionalitas kode yang telah ditulis terhadap instruksi modul. AI menemukan ketidaksesuaian kecil (seperti tipe data `Grade` yang harusnya `float64` alih-alih `int`) dan menyarankan penambahan fungsi pembanding untuk *pass by value*.
- **Manajemen Repositori:** Meminta saran praktik terbaik (*best practices*) untuk menyusun struktur folder proyek Go agar terhindar dari bentrok *package main* ketika ada banyak file tugas dalam satu repositori.

### 📚 Referensi Eksternal
Selain modul dan AI, referensi berikut juga digunakan untuk mendalami pengerjaan tugas:
1. **Dokumentasi Resmi Golang (Package `fmt`)**: Digunakan untuk mempelajari *formatting verbs* saat melakukan print output ke terminal.
   - Tautan: [https://pkg.go.dev/fmt](https://pkg.go.dev/fmt)
2. **StackOverflow**: Digunakan untuk membaca diskusi programmer lain mengenai konsep *Pass by Reference* vs *Pass by Value* pada bahasa Go.
   - Tautan: [https://stackoverflow.com/questions/47296325/passing-by-reference-and-value-in-go-to-functions](https://stackoverflow.com/questions/47296325/passing-by-reference-and-value-in-go-to-functions)

---

## Pertemuan 2 (TM 2)

### 🤖 Bantuan AI
Dalam pengerjaan TM 2, AI digunakan sebagai asisten untuk menjelaskan alur kerja API, melakukan *code review*, dan membantu menyusun draf laporan, dengan rincian sebagai berikut:
- **Penjelasan Modul & Konsep:** Membantu merangkum ulang materi modul, memecah alur *request* dengan diagram visual, dan mempertegas poin penting seperti alasan penggunaan tipe data *pointer* pada metode PATCH dibanding PUT.
- **Eksplorasi Konsep Tambahan:** Menjelaskan secara teoritis konsep-konsep seputar HTTP dan REST API yang tidak tercantum di modul (seperti *idempotency*, perbedaan JSON Patch vs Merge Patch, *integer division* di Go, serta cara memvalidasi data unik atau duplikat pada *array in-memory* sebelum beralih ke database).
- **Perancangan Tugas Mandiri:** Memandu alur penyelesaian tugas API Students dari awal hingga penyusunan *query string* untuk paginasi dan filter logika tambahan (seperti batas nilai *grade*).
- **Code Review & Debugging:** AI digunakan untuk meninjau dan mengevaluasi kode proyek `api-students`. AI berhasil mendeteksi dan mengusulkan perbaikan untuk beberapa *bug* krusial pada logika filter (*typo* parameter *query* dan kesalahan *parsing* boolean), logika *slicing* pada metode DELETE, hingga pencegahan duplikasi NIM (409 Conflict).

### 📚 Referensi Eksternal
Selain modul dan bantuan AI, referensi eksternal berikut juga digunakan untuk mendukung pengerjaan tugas:
1. **StackOverflow**: Digunakan sebagai referensi bacaan diskusi tentang pembuatan parameter unik dari *struct* di Go.
   - Tautan: [https://stackoverflow.com/questions/48253423/unique-hash-from-struct](https://stackoverflow.com/questions/48253423/unique-hash-from-struct)
2. **Dokumentasi Resmi Go Fiber**: Digunakan untuk membaca panduan dan dokumentasi framework web Fiber versi 2.
   - Tautan: [https://pkg.go.dev/github.com/gofiber/fiber/v2@v2.52.15](https://pkg.go.dev/github.com/gofiber/fiber/v2@v2.52.15)
3. **MDN Web Docs (HTTP CORS)**: Digunakan untuk mempelajari landasan teori mengenai *Cross-Origin Resource Sharing* (CORS) dalam protokol HTTP.
   - Tautan: [https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS)
4. **MDN Web Docs (HTTP General)**: Digunakan untuk mendalami standar komunikasi web berbasis protokol HTTP.
   - Tautan: [https://developer.mozilla.org/en-US/docs/Web/HTTP](https://developer.mozilla.org/en-US/docs/Web/HTTP)

---
## Pertemuan 3 (TM 3)

### 🤖 Bantuan AI
Dalam pengerjaan TM 3, AI digunakan sebagai tutor arsitektur, pemandu migrasi kode, serta asisten pengujian, dengan rincian sebagai berikut:
- **Penjelasan Modul & Konsep Database:** Membantu membedah isi modul mengenai pola *Repository* dan integrasi PostgreSQL. AI memberikan penjelasan mendalam (*deep dive*) tentang properti ACID, perbedaan `pgx` dengan `database/sql`, mekanisme *parameterized query* untuk mencegah SQL *Injection*, serta bahaya performa pada perintah `OFFSET`.
- **Eksplorasi Konsep Tambahan:** Menjelaskan secara teoritis fungsi klausul `WHERE 1 = 1` dalam penyusunan query dinamis, cara kerja pemetaan (*scanning*) data hasil query ke dalam *struct* dengan pointer, dan pemahaman terkait pembongkaran error berlapis (*error wrapping*) menggunakan fungsi `errors.As()`.
- **Perancangan & Migrasi Tugas Mandiri:** Memberikan panduan *step-by-step* untuk memigrasikan API `api-students` dari penyimpanan memori (TM 2) ke database relasional. Panduan ini mencakup perombakan *handler*, pembuatan interface *repository*, penyusunan migrasi SQL yang aman (*Unique Index*), serta injeksi dependensi.
- **Code Review & Debugging:** AI digunakan untuk meninjau proyek akhir TM 3. AI berhasil mendeteksi dan mengusulkan perbaikan terhadap *bug* pemetaan target tabel saat operasi DELETE, dan perbaikan penulisan target pointer (`*req.NIM`) saat operasi validasi di dalam metode PATCH.
- **Pengujian (Testing):** Membantu menyusun rancangan perintah pengujian via cURL untuk 21 skenario, serta menyusun panduan penggunaan *Insomnia REST Client* untuk skenario duplikasi data, filter SQL Injection, hingga kasus database *down* (500/503).

### 📚 Referensi Eksternal
Selain modul dan bantuan AI, referensi eksternal berikut juga digunakan untuk mendukung pengerjaan tugas:
1. **Dokumentasi pgxpool (v5)**: Digunakan untuk mempelajari cara mengonfigurasi dan menggunakan *connection pool* pada driver PostgreSQL `pgx` versi 5.
   - Tautan: [https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool)
2. **Dokumentasi pgx**: Digunakan untuk membaca panduan standar terkait implementasi koneksi Go dan PostgreSQL.
   - Tautan: [https://pkg.go.dev/github.com/jackc/pgx](https://pkg.go.dev/github.com/jackc/pgx)
3. **Daftar Kode Error PostgreSQL**: Digunakan untuk mempelajari kode-kode *error* bawaan PostgreSQL, khususnya terkait pelanggaran konstrain unik (`23505`).
   - Tautan: [https://www.postgresql.org/docs/current/errcodes-appendix.html](https://www.postgresql.org/docs/current/errcodes-appendix.html)
4. **RFC 9110 (HTTP Semantics)**: Digunakan sebagai pedoman standar terbaru terkait definisi dan semantik kode status HTTP, khususnya untuk memvalidasi penggunaan status 503 (*Service Unavailable*).
   - Tautan: [https://www.rfc-editor.org/info/rfc9110/](https://www.rfc-editor.org/info/rfc9110/)

---

## Pertemuan 4 (TM 4)

### 🤖 Bantuan AI
Dalam pengerjaan TM 4, AI digunakan sebagai tutor arsitektur, *code reviewer*, serta asisten refactoring, dengan rincian sebagai berikut:
- **Penjelasan Modul & Konsep Arsitektur:** Membantu membedah isi modul mengenai *Clean Architecture*. AI memberikan penjelasan mendalam tentang empat *layer* arsitektur, *Dependency Rule*, penerapan *Dependency Inversion* menggunakan *interface*, serta kaitannya dengan prinsip desain *SOLID*.
- **Eksplorasi Konsep Tambahan:** Menjelaskan secara teoritis konsep-konsep di luar modul seperti penggunaan `context.WithTimeout` untuk mencegah *resource exhaustion*, mekanisme *Graceful Shutdown*, konsep *Log Rotation* menggunakan Lumberjack, pola *Closure* pada *middleware*, hingga *Sentinel Errors* menggunakan `errors.Is()`.
- **Pemahaman Alur Kerja Baru:** Membantu memetakan ulang alur permintaan (*request flow*) dari API yang tadinya *flat* menjadi terstruktur sesuai pembagian tanggung jawab (*routing* → *middleware* → *controller* → *business rules* → *repository*).
- **Code Review & Audit Proyek:** AI digunakan untuk mengaudit dan mengevaluasi proyek `latihan-fiber` dan `api-students` terhadap standar modul. AI berhasil mendeteksi *bug* tersembunyi serta menemukan *file* lama yang belum terhapus.
- **Refactoring & Perbaikan Kode:** Membantu menerapkan solusi untuk menyelesaikan seluruh *bug* hasil audit, memandu perbaikan logika pada fungsi `ApplyPatch`, serta memperbarui dokumentasi proyek di `README.md` dan `.gitignore` agar mencerminkan struktur yang baru.

### 📚 Referensi Eksternal
Selain modul dan bantuan AI, referensi eksternal berikut juga digunakan untuk mendukung pengerjaan tugas:
1. **Dokumentasi Resmi Golang (Package `log/slog`)**: Digunakan untuk mempelajari cara implementasi *structured logging* bawaan Go.
   - Tautan: [https://pkg.go.dev/log/slog](https://pkg.go.dev/log/slog)
2. **Dokumentasi Pustaka Lumberjack (v2)**: Digunakan untuk mempelajari konfigurasi rotasi *file log* otomatis (berdasarkan ukuran, usia, dan batas pencadangan).
   - Tautan: [https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2](https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2)
3. **Dokumentasi Resmi Golang (Package `path/filepath`)**: Digunakan untuk memahami manipulasi lintasan (*path*) direktori secara aman lintas sistem operasi saat mengonfigurasi lokasi *file log*.
   - Tautan: [https://pkg.go.dev/path/filepath](https://pkg.go.dev/path/filepath)
4. **Dokumentasi Resmi Golang (Package `os`)**: Digunakan untuk memahami interaksi langsung dengan sistem operasi, seperti pembuatan folder (`MkdirAll`), penanganan sinyal terminasi (`os.Signal`), dan keluar program (`os.Exit`).
   - Tautan: [https://pkg.go.dev/os](https://pkg.go.dev/os)
5. **Dokumentasi Resmi Golang (Package `io`)**: Digunakan untuk mempelajari antarmuka I/O primitif, khususnya pola `MultiWriter` yang memungkinkan output log ditulis serentak ke layar (terminal) dan *file*.
   - Tautan: [https://pkg.go.dev/io](https://pkg.go.dev/io)

---

## Pertemuan 5 (TM 5)

### 🤖 Bantuan AI
Dalam pengerjaan TM 5, AI digunakan sebagai tutor arsitektur keamanan, *code reviewer*, serta asisten pengujian, dengan rincian sebagai berikut:
- **Penjelasan Modul & Konsep Keamanan:** Membantu membedah isi modul mengenai autentikasi vs otorisasi, cara kerja *hashing* password menggunakan Bcrypt (salt dan cost factor), struktur token JWT, serta pertimbangan penyimpanan token (*localStorage* vs *httpOnly cookie*).
- **Eksplorasi Konsep Tambahan:** Menjelaskan secara teoritis topik keamanan lanjutan di luar modul, seperti perbandingan Bcrypt vs Argon2id, mitigasi serangan *Algorithm Confusion* (`alg: none`), *Timing Attack*, *Mass Assignment*, dan mekanisme rotasi *refresh token*.
- **Code Review & Debugging:** AI digunakan untuk mengevaluasi kode dan menemukan kendala, seperti perbaikan sintaks pada migrasi SQL, *case sensitivity* pada unit test pesan error validasi password, serta memastikan fungsi validasi bersifat *pure function*.

### 📚 Referensi Eksternal
Selain modul dan bantuan AI, referensi eksternal berikut juga digunakan untuk mendukung pengerjaan tugas:
1. **Dokumentasi golang-jwt (v5)**: Digunakan untuk mempelajari pembuatan, penandatanganan, dan parsing JSON Web Token (JWT).
   - Tautan: [https://pkg.go.dev/github.com/golang-jwt/jwt/v5](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)
2. **Dokumentasi golang-jwt (README)**: Digunakan sebagai panduan implementasi dan contoh penggunaan pustaka `jwt/v5`.
   - Tautan: [https://pkg.go.dev/github.com/golang-jwt/jwt/v5#section-readme](https://pkg.go.dev/github.com/golang-jwt/jwt/v5#section-readme)
3. **Dokumentasi bcrypt**: Digunakan untuk memahami fungsi hash password adaptif `GenerateFromPassword` dan `CompareHashAndPassword`.
   - Tautan: [https://pkg.go.dev/golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
4. **Dokumentasi crypto/sha256**: Digunakan untuk implementasi hashing SHA-256 pada penyimpanan string refresh token.
   - Tautan: [https://pkg.go.dev/crypto/sha256](https://pkg.go.dev/crypto/sha256)
5. **Dokumentasi crypto/rand**: Digunakan untuk menghasilkan byte acak berbasis kriptografi yang aman saat membuat refresh token.
   - Tautan: [https://pkg.go.dev/crypto/rand](https://pkg.go.dev/crypto/rand)
6. **Dokumentasi encoding/hex**: Digunakan untuk mengubah representasi byte acak dan hash menjadi string heksadesimal.
   - Tautan: [https://pkg.go.dev/encoding/hex](https://pkg.go.dev/encoding/hex)
7. **Dokumentasi errors**: Digunakan untuk mencocokkan jenis error menggunakan `errors.Is()` pada penanganan token kedaluwarsa vs tidak valid.
   - Tautan: [https://pkg.go.dev/errors](https://pkg.go.dev/errors)
8. **Dokumentasi crypto**: Digunakan sebagai referensi standar antarmuka kriptografi bawaan Go.
   - Tautan: [https://pkg.go.dev/crypto](https://pkg.go.dev/crypto)
9. **Dokumentasi builtin (make)**: Digunakan untuk memahami alokasi slice byte penampung data acak.
   - Tautan: [https://pkg.go.dev/builtin#make](https://pkg.go.dev/builtin#make)
10. **Dokumentasi builtin**: Digunakan sebagai referensi tipe data dan fungsi dasar bawaan Go.
    - Tautan: [https://pkg.go.dev/builtin](https://pkg.go.dev/builtin)

