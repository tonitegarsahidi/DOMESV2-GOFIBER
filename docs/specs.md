# Backend System Specifications — DOMES V2

Dokumen ini menjelaskan spesifikasi teknis, arsitektur, database, modul, serta kontrak API pada bagian Backend proyek **DOMES V2** berbasis Go Fiber.

---

## 🛠️ Technology Stack

* **Programming Language**: Go 1.21
* **Web Framework**: [GoFiber v2](https://gofiber.io/) (`github.com/gofiber/fiber/v2`) - web framework berbasis fasthttp dengan performa tinggi.
* **Database & ORM**: MySQL dengan [GORM](https://gorm.io/) (`gorm.io/gorm`, `gorm.io/driver/mysql`)
* **Cache**: Redis (`github.com/redis/go-redis/v9`) (opsional, dikonfigurasi melalui `.env`)
* **Logging**: Zap Logger (`go.uber.org/zap`)
* **Security & Auth**: 
  * JWT (`github.com/golang-jwt/jwt/v5`)
  * Bcrypt untuk enkripsi kata sandi (`golang.org/x/crypto/bcrypt`)
  * Google reCAPTCHA v2 (`pkg/captcha`)
* **Environment Configuration**: godotenv (`github.com/joho/godotenv`)

---

## 🏗️ Architecture Pattern (Controller-Service-Repository)

Backend menggunakan arsitektur berlapis untuk menjaga pemisahan tanggung jawab (*separation of concerns*):
1. **Entrypoint** (`cmd/main.go`): Menginisialisasi koneksi DB, Logger, Redis, mengaktifkan Auto-Migration & Seeders, lalu menjalankan server Fiber.
2. **Routes** (`routes/routes.go`): Defini rute HTTP dan penempatan middleware yang sesuai (seperti JWT auth middleware).
3. **Controller Layer** (`internal/controller/`): Menangani parsing HTTP request (JSON body, query params), validasi dasar, dan formatting response.
4. **Service Layer** (`internal/service/`): Berisi seluruh logika bisnis aplikasi, koordinasi transaksi, serta integrasi email.
5. **Repository Layer** (`internal/repository/`): Menangani query langsung ke database menggunakan GORM ORM.
6. **Model Layer** (`internal/model/`): Definisi skema tabel GORM dan tipe data terstruktur.

---

## 🗄️ Database Architecture (V2 Schema)

Semua tabel database telah ditingkatkan ke struktur V2 dengan spesifikasi sebagai berikut:
* **Table Prefix**: Nama tabel berawalan `V2` (misal: `V2Documents`, `V2AdminEmails`, `v2_document_sdgs`, dll.).
* **Primary Key**: Menggunakan UUID v4 string (bukan auto-increment integer) untuk mencegah enumerasi ID dokumen secara publik.
* **Timestamps**: Kolom default `CreatedAt` dan `UpdatedAt` bertipe timestamp.
* **Audit Fields**: Kolom `CreatedBy` dan `UpdatedBy` menyimpan ID pengguna yang membuat atau mengubah record data.
* **Soft Delete**: Kolom `DeletedAt` disediakan untuk penanganan soft delete data.
* **Auto-Migrations**: Dijalankan secara otomatis saat server berjalan (kecuali tabel `Users` yang membutuhkan environment `RUN_USER_MIGRATION=true`).

### Relasi Dokumen & Metadata Master:
* Relasi Many-to-Many terwujud melalui tabel perantara seperti `v2_document_sdgs`, `v2_document_sectors`, `v2_document_lnobs`, dan `v2_document_agencies`.
* Skrip migrasi data lama (`go run cmd/migrate_data/main.go`) memetakan data legacy dari tabel `Tabledatas` ke skema V2 secara terstruktur.

---

## 📡 API Endpoints Summary

### Rute Publik (Public Routes)
* **`GET /api/v2/health-check`**: Pengecekan status server, koneksi database, dan Redis.
* **`GET /api/v2/documents`**: Menampilkan daftar dokumen terfilter dengan pagination, pencarian teks penuh, filter SDG, instansi, sektor, LNOB, dan bahasa.
* **`GET /api/v2/documents/:id`**: Detail dokumen berdasarkan ID (UUID) atau Slug.
* **`GET /api/v2/master/languages`**: List data master bahasa yang tersedia.
* **`GET /api/v2/master/sdgs`**: List data master Sustainable Development Goals (SDGs).
* **`GET /api/v2/master/agencies`**: List data master UN Agencies.
* **`GET /api/v2/master/sectors`**: List data master thematic sectors.
* **`GET /api/v2/master/lnobs`**: List data master Leave No One Behind (LNOB) categories.

### Rute Otentikasi (Auth Routes)
* **`POST /api/v2/auth/login`**: Login admin/kontributor dengan validasi captcha.
* **`POST /api/v2/auth/register`**: Registrasi akun kontributor baru.
* **`POST /api/v2/auth/forgot-password`**: Permintaan reset password.

### Rute Admin & Kontributor (Protected Routes - `/api/v2/cms/*` dan `/api/v2/submissions/*`)
*Diperlukan Header `Authorization: Bearer <JWT_TOKEN>`*
* **`GET /api/v2/cms/dashboard-stats`**: Statistik ringkasan pengajuan dokumen.
* **`GET /api/v2/submissions`**: Daftar antrean pengajuan dokumen.
* **`POST /api/v2/submissions`**: Mengajukan dokumen baru (multi-step data).
* **`PUT /api/v2/submissions/:id`**: Memperbarui dokumen yang diajukan.
* **`POST /api/v2/submissions/:id/draft`**: Menyimpan draft dokumen per step.
* **`PUT /api/v2/submissions/:id/publish`**: Mempublikasikan dokumen.
* **`PUT /api/v2/submissions/:id/unpublish`**: Membatalkan publikasi dokumen.
* **`POST /api/v2/upload`**: Upload file PDF dokumen (disimpan di folder `routes/uploads/` dengan nama berkas berbasis UUID).
* **`GET/POST/PUT/DELETE /api/v2/cms/master/*`**: CRUD pengelolaan data master.
