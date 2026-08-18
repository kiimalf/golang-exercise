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
