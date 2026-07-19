# Backend Changelog — DOMES V2

Dokumen ini mencatat riwayat perubahan, pembaruan fitur, perbaikan bug, dan migrasi database yang dilakukan pada bagian Backend **DOMES V2** (Go Fiber).

---

## [1.0.2] - 2026-07-12

### Added
- **Intermediate File Upload Service**: Memperbarui backend untuk mendelegasikan upload file ke `FileUploadService` perantara, yang mengontrol tipe file (`main`, `cover`, `additional`), melakukan evaluasi metadata/ukuran file, dan menformat file sebagai `TYPE-uuid.extension` ke folder `./public/upload`.
- **Create Submission Draft Flow**: Mengubah `CreateSubmission` pada `DocumentService` agar mengizinkan pembuatan dokumen tanpa judul (defaulting ke `"Draft Submission"`) dan menetapkan status default dokumen menjadi `"draft"` (sebelumnya langsung `"pending_review"`), memfasilitasi pembuatan draft instan sejak Step 1.

### Changed
- **Title Validation Relaxation in Updates**: Menghilangkan validasi wajib judul (`Title`) pada `UpdateSubmission` saat judul dikosongkan (sistem akan mempertahankan judul draft yang ada), sehingga user bebas bolak-balik mengubah data antarlangkah wizard tanpa memicu error validasi judul.

---

## [1.0.1] - 2026-07-12

### Added
- **Update Submission Endpoint (PUT)**: Menambahkan endpoint `PUT /api/v2/submissions/:id` pada layer route, controller, dan service. Endpoint ini memproses data `SubmissionRequest` lengkap, memverifikasi kepemilikan dokumen (AuthorID == userID atau role administrator), memperbarui metadata dasar, menyelaraskan many-to-many associations (SDGs, Sectors, LNOBs) menggunakan GORM association replacement, dan memperbarui slug secara dinamis jika judul dokumen diubah.

### Changed
- **Database Connection Retry with Backoff**: Mengganti koneksi database satu kali (*fire-and-forget*) menjadi mekanisme retry otomatis dengan exponential backoff (5 percobaan: 1s → 2s → 4s → 8s → 16s). Server kini menunggu database siap sebelum menerima traffic, alih-alih berjalan dengan `DB = nil` yang menyebabkan nil pointer panic pada setiap request. Jika database tetap tidak tersedia setelah semua retry, server akan exit dengan pesan error yang jelas. Logging juga dimigrasikan dari `log.Printf` ke Zap Logger sesuai code style guidelines.

---

## [1.0.0] - 2026-07-11

### Added
- **Dynamic Master Metadata API**: Endpoint `/api/v2/master/languages` untuk menyediakan daftar bahasa dinamis bagi frontend.
- **Reference Table Seeders**: Seeder otomatis saat server dijalankan untuk mengisi data dasar SDGs, Agencies, Sectors, dan LNOBs.

### Changed
- **Response Format Standardization**: Merapikan response HTTP error agar mengembalikan format standard yang seragam (`{ "success": false, "message": "..." }`) saat request gagal, menghindari bocornya stack trace internal ke pengguna umum.
- **Seeder & Migration Safety**: Menambahkan parameter keamanan `RUN_USER_MIGRATION` agar migrasi tabel pengguna (`Users`) dilewati demi mencegah terhapusnya akun operasional secara tidak sengaja saat menjalankan deployment.
- **Related Documents Query Limit**: Mengubah limit pencarian dokumen terkait (*related documents*) pada level database/repository GORM dari 3 menjadi 4 untuk memfasilitasi kebutuhan UI frontend.

### Fixed
- **Testing Database Conflict**: Memperbaiki issue di mana eksekusi unit test (`go test ./...`) secara tidak sengaja menghapus tabel database utama. Migrasi test kini diisolasi sepenuhnya.
- **Document Detail Lookup**: Memperbaiki pencarian dokumen di endpoint `GET /api/v2/documents/:id` agar mendukung query pencarian baik berdasarkan UUID (ID dokumen) maupun Slug.

---

## [0.9.0] - 2026-07-03

### Added
- **Legacy Migrator**: Membuat command script manual di `cmd/migrate_data/main.go` untuk memetakan data lama dari tabel `Tabledatas` ke skema tabel `V2Documents` beserta seluruh relasi Many-to-Many secara teratur.
- **CMS Stats Endpoint**: Menambahkan endpoint `/api/v2/cms/dashboard-stats` untuk kalkulasi analitik ringkasan dokumen per status persetujuan.

### Changed
- **Migration Schema V2**: Migrasi penuh skema tabel database lama dengan awalan `V2` (misalnya `V2Documents`, `V2AdminEmails`) serta merubah tipe primary key menjadi string UUID v4.

---

## [0.8.0] - 2026-06-25

### Added
- **Multi-step Submission REST Handler**: Endpoint `POST /api/v2/cms/submissions` untuk menangani input multi-step pengajuan dokumen dari kontributor.
- **File Upload Handler**: Integrasi upload berkas PDF ke subfolder `routes/uploads` dengan konversi nama berkas menjadi UUID otomatis demi privasi data berkas.
- **Auth Guard Middleware**: Implementasi JWT Validation Middleware untuk memproteksi endpoint admin `/api/v2/cms/*`.
