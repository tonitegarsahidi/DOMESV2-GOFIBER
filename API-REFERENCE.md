# API Reference - DOMESv2

Dokumentasi lengkap API endpoints untuk DOMESv2 Backend.

## Base URL

| Environment | URL |
|-------------|-----|
| Development | `http://localhost:3000` |
| Production  | `https://domesv2.yourdomain.com` |

## Authentication

API menggunakan **JWT (JSON Web Token)** untuk autentikasi.
- Token didapat dari endpoint `/api/v2/auth/login` atau `/api/v2/auth/register`
- Kirim token di header: `Authorization: Bearer <token>`
- Token berlaku 24 jam (default)

---

## Ringkasan Endpoints

### Public Endpoints (Tanpa Auth)
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/` | Informasi dasar API |
| GET | `/api/v2/health-check` | Cek status kesehatan sistem |
| POST | `/api/v2/auth/register` | Pendaftaran user baru |
| POST | `/api/v2/auth/login` | Login user |
| POST | `/api/v2/auth/forgot-password` | Pengiriman email reset password |
| POST | `/api/v2/auth/reset-password` | Reset password dengan token |
| GET | `/api/v2/reference/:type` | Mengambil data referensi (agencies, sdgs, sectors, dll.) |
| GET | `/api/v2/documents` | List dokumen terbitan publik |
| GET | `/api/v2/documents/search` | Pencarian teks bebas dokumen |
| GET | `/api/v2/documents/:id` | Detail dokumen berdasarkan ID atau slug |
| GET | `/api/v2/documents/:id/related` | Rekomendasi dokumen terkait |
| GET | `/api/v2/documents/:id/download` | Unduh file dokumen (dan tracking download) |
| GET | `/api/v2/stats` | Statistik agregat platform |
| GET | `/api/v2/analytics/:type` | Grafik analitik publik (overview, by-sdg, dll.) |
| POST | `/api/v2/reports` | Mengirim laporan broken link |

### Protected Endpoints (Memerlukan JWT Bearer Token)
| Method | Endpoint | Role | Deskripsi |
|--------|----------|------|-----------|
| GET | `/api/v2/user/me` | User & Admin | Profil user login |
| PUT | `/api/v2/user/profile` | User & Admin | Edit profil diri |
| PUT | `/api/v2/user/password` | User & Admin | Ganti password |
| GET | `/api/v2/user/notifications` | User & Admin | Ambil preferensi notifikasi |
| PUT | `/api/v2/user/notifications` | User & Admin | Edit preferensi notifikasi |
| GET | `/api/v2/admin/emails` | Admin | List email whitelist admin |
| POST | `/api/v2/admin/emails` | Admin | Tambah email ke whitelist admin |
| DELETE | `/api/v2/admin/emails/:email` | Admin | Hapus email dari whitelist |
| GET | `/api/v2/cms/dashboard` | User & Admin | Ringkasan dashboard editor |
| GET | `/api/v2/cms/activity` | User & Admin | Aktivitas terbaru editor |
| GET | `/api/v2/submissions` | User & Admin | List dokumen pengajuan editor |
| POST | `/api/v2/submissions` | User & Admin | Submit final dokumen (Step 4) |
| POST | `/api/v2/submissions/:id/draft` | User & Admin | Simpan draf dokumen (Step 1-3) |
| DELETE | `/api/v2/submissions/:id` | User & Admin | Hapus dokumen pengajuan |
| PUT | `/api/v2/submissions/:id/publish` | User & Admin | Terbitkan dokumen ke publik |
| PUT | `/api/v2/submissions/:id/unpublish` | User & Admin | Tarik dokumen dari publik |
| GET | `/api/v2/reports` | User & Admin | List laporan broken link |
| PUT | `/api/v2/reports/:id/status` | User & Admin | Update status laporan link |
| GET | `/api/v2/analytics/summary` | User & Admin | Ringkasan analitik internal |
| GET | `/api/v2/analytics/top-downloads` | User & Admin | Top dokumen terunduh |
| GET | `/api/v2/analytics/top-views` | User & Admin | Top dokumen dilihat |
| POST | `/api/v2/upload` | User & Admin | Upload file dokumen/media |
| POST | `/api/v2/upload/url-validate` | User & Admin | Validasi URL eksternal |
| POST | `/api/v2/upload/avatar` | User & Admin | Upload foto profil/avatar |
| GET | `/api/v2/users` | Admin | List user pengelola |
| POST | `/api/v2/users` | Admin | Tambah user pengelola manual |
| PUT | `/api/v2/users/:id` | Admin | Edit user pengelola |
| DELETE | `/api/v2/users/:id` | Admin | Hapus user pengelola |
| GET | `/api/v2/cms/reference/:type` | User & Admin | List reference data CMS |
| POST | `/api/v2/cms/reference/:type` | Admin | Tambah reference data CMS |
| PUT | `/api/v2/cms/reference/:type/:code` | Admin | Edit reference data CMS |
| DELETE | `/api/v2/cms/reference/:type/:code` | Admin | Hapus reference data CMS |

---

## Detail Endpoints

### 1. Authentication & Profiles

#### POST `/api/v2/auth/register` (Public)
Mendaftarkan user baru. Role diatur otomatis berdasarkan daftar whitelist admin.
* **Request Body:**
  ```json
  {
    "first_name": "Toni",
    "last_name": "Tegar",
    "position": "Staff",
    "organization": "UNDP",
    "phone_number": "+62812345678",
    "email": "tonitegarsahidi@gmail.com",
    "password": "rahasiaku123",
    "confirm_password": "rahasiaku123",
    "captcha": "optional-recaptcha-token"
  }
  ```
* **Response 201:** Sukses mendaftar, mengembalikan JWT token & data user.

#### POST `/api/v2/auth/login` (Public)
Login untuk mendapatkan token akses.
* **Request Body:**
  ```json
  {
    "email": "tonitegarsahidi@gmail.com",
    "password": "rahasiaku123",
    "captcha": "optional-recaptcha-token"
  }
  ```
* **Response 200:**
  ```json
  {
    "success": true,
    "message": "Login successful",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "user": { "id": 1, "email": "tonitegarsahidi@gmail.com", "role": "administrator", ... }
    }
  }
  ```

#### GET `/api/v2/user/me` (Protected)
Mendapatkan data lengkap user yang sedang login beserta preferensinya.
* **Response 200:** Mengembalikan objek user lengkap.

#### PUT `/api/v2/user/profile` (Protected)
Mengubah nama depan, nama belakang, telepon, organisasi, atau jabatan.
* **Request Body:**
  ```json
  {
    "first_name": "Toni",
    "last_name": "Sahidi",
    "position": "Senior Staff",
    "organization": "UNDP",
    "phone_number": "+6289999999"
  }
  ```

#### PUT `/api/v2/user/password` (Protected)
Mengganti password akun aktif.
* **Request Body:**
  ```json
  {
    "current_password": "rahasiaku123",
    "new_password": "passwordbaru123",
    "confirm_password": "passwordbaru123"
  }
  ```

#### GET & PUT `/api/v2/user/notifications` (Protected)
Mengambil atau mengupdate preferensi notifikasi sistem.
* **Request Body (PUT):**
  ```json
  {
    "document_approvals": true,
    "broken_link_reports": false,
    "system_updates": true,
    "email_notifications": true
  }
  ```

---

### 2. Admin Whitelist Settings (Protected - Admin Only)

#### GET `/api/v2/admin/emails`
Mengambil daftar whitelist admin email.
* **Response 200:** Array objek email whitelist.

#### POST `/api/v2/admin/emails`
Menambahkan email baru agar saat mendaftar otomatis menjadi Administrator.
* **Request Body:**
  ```json
  { "email": "new-admin@un.org" }
  ```

#### DELETE `/api/v2/admin/emails/:email`
Menghapus email dari daftar whitelist.

---

### 3. Reference Data (Public)

#### GET `/api/v2/reference/:type`
* Parameter `:type` yang valid: `agencies`, `sdgs`, `sectors`, `languages`, `joint-programmes`, `lnobs`, `non-un-partners`, `organizations`.
* **Response 200:** List objek data referensi (berisi id, code, name, icon/color jika ada).

---

### 4. Public Documents Discovery (Public)

#### GET `/api/v2/documents`
Mendapatkan semua dokumen yang berstatus `published` dengan pagination & filter.
* **Query Parameters:**
  * `page` (default 1)
  * `limit` (default 10)
  * `agency` (code agensi)
  * `sdg` (code SDG)
  * `sector` (code sektor)
  * `language` (code bahasa)
  * `sort` (`newest`, `oldest`, `downloads`, `views`)

#### GET `/api/v2/documents/search`
Pencarian teks bebas pada dokumen.
* **Query Parameters:** `q` (kata kunci pencarian), `sort` (`relevance`, `newest`, dsb.)

#### GET `/api/v2/documents/:id`
Mencari dokumen berdasarkan ID UUID v4 atau Slug teks unik.

#### GET `/api/v2/documents/:id/related`
Mendapatkan rekomendasi dokumen lain yang memiliki irisan SDG atau sektor.

#### GET `/api/v2/documents/:id/download`
Meningkatkan counter downloads dokumen dan mengembalikan link unduhan.

---

### 5. Broken Link Reporting (Public & Protected)

#### POST `/api/v2/reports` (Public)
Mengajukan laporan link PDF dokumen yang rusak/404.
* **Request Body:**
  ```json
  {
    "document_id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
    "reporter_name": "John Doe",
    "reporter_email": "johndoe@example.com",
    "details": "Tautan download PDF mengarah ke halaman kosong."
  }
  ```

#### GET `/api/v2/reports` (Protected)
Mengambil daftar laporan yang diajukan. Query param: `status` (`all`, `pending`, `in_progress`, `resolved`).

#### PUT `/api/v2/reports/:id/status` (Protected)
Memperbarui status penanganan laporan broken link.
* **Request Body:**
  ```json
  { "status": "resolved" } // Pilihan: pending, in_progress, resolved
  ```

---

### 6. CMS & Submissions Management (Protected)

#### GET `/api/v2/cms/dashboard` & `/api/v2/cms/activity`
Statistik performa CMS internal editor dan log aktivitas riwayat aksi terbaru.

#### POST `/api/v2/submissions/:id/draft`
Menyimpan progres langkah wizard submission (Step 1-3).
* **Request Body:**
  ```json
  {
    "step": 2,
    "data": {
      "title": "Draf Dokumen Baru",
      "short_summary": "Summary draf...",
      "focal_point_name": "Budi"
    }
  }
  ```

#### POST `/api/v2/submissions`
Mengirimkan data dokumen lengkap secara final (Step 4).
* **Request Body:**
  ```json
  {
    "title": "Dokumen Laporan Final",
    "short_description": "Deskripsi dokumen...",
    "abstract": "Abstract dokumen...",
    "detailed_summary": "Summary lengkap HTML...",
    "date_of_publication": "2024-06-24",
    "total_pages": 88,
    "language": "English",
    "publication_status": "Published",
    "tags": ["economy", "sustainability"],
    "file_url": "/uploads/example.pdf",
    "file_size": "2.4 MB",
    "cover_image_url": "/uploads/cover.jpg",
    "agency": "UNDP",
    "focal_point": { "name": "Budi", "email": "budi@un.org", "phone": "0812...", "department": "R&D" },
    "sdgs": ["GOAL 1", "GOAL 8"],
    "sectors": ["economic-development"],
    "lnob_groups": ["disabilities"],
    "joint_programme": "adlight"
  }
  ```

#### PUT `/api/v2/submissions/:id/publish` & `/api/v2/submissions/:id/unpublish`
Mengubah visibilitas dokumen (memublikasikan ke portal publik atau menyembunyikan).

---

### 7. File Upload & Validation (Protected)

#### POST `/api/v2/upload`
Mengunggah file. Menggunakan parser multipart form-data.
* **Request Form:**
  * `file`: File media (PDF, Word, JPG, PNG)
  * `type`: `document` (untuk PDF/Word) atau `cover` (untuk image)
* **Response 201:** `{"success": true, "url": "/uploads/random-uuid.pdf", "size": "1.2 MB"}`

#### POST `/api/v2/upload/url-validate`
Validasi URL eksternal apakah merespons dengan HTTP status 200 OK.
* **Request Body:** `{ "url": "https://active-link.com/document.pdf" }`

#### POST `/api/v2/upload/avatar`
Upload foto profil diri (avatar) user aktif. Form key: `avatar` (image).

---

### 8. CMS User Management (Protected - Admin Only)

#### GET, POST, PUT, DELETE pada `/api/v2/users`
Fungsi manajemen akun pengelola (CRUD) oleh Admin.
* **Request Body (POST - Create User):**
  ```json
  {
    "first_name": "Jane",
    "last_name": "Doe",
    "email": "janedoe@un.org",
    "password": "password123",
    "confirm_password": "password123",
    "organization": "FAO",
    "position": "Editor Staff",
    "role": "editor",
    "status": "active"
  }
  ```

---

### 9. CMS Reference Data Management (Protected)

#### GET `/api/v2/cms/reference/:type` (Protected)
Mengambil daftar reference data untuk `:type` tertentu (seperti `agencies`, `sdgs`, `sectors`, `languages`, `joint-programmes`, `lnobs`, `non-un-partners`, `organizations`).
* **Response 200:** Array objek reference.

#### POST `/api/v2/cms/reference/:type` (Protected - Admin Only)
Menambahkan item referensi baru ke database.
* **Request Body:**
  ```json
  {
    "code": "test-sector-xyz",
    "name": "Test Sector XYZ",
    "logo_url": "optional-logo-url-for-agencies",
    "icon": "optional-icon-for-sdgs",
    "color": "optional-color-for-sdgs"
  }
  ```
* **Response 201:** Objek referensi yang berhasil dibuat.

#### PUT `/api/v2/cms/reference/:type/:code` (Protected - Admin Only)
Memperbarui informasi nama atau metadata item referensi berdasarkan kode primary key.
* **Request Body:**
  ```json
  {
    "name": "Updated Test Sector XYZ",
    "logo_url": "updated-logo-url"
  }
  ```
* **Response 200:** Objek referensi ter-update.

#### DELETE `/api/v2/cms/reference/:type/:code` (Protected - Admin Only)
Menghapus item referensi dari database.
* **Response 200:** Sukses menghapus referensi.

---

## Error Codes

| Kode | HTTP Status | Deskripsi |
|------|-------------|-----------|
| `INVALID_REQUEST_BODY` | 400 | Format request body salah |
| `INVALID_RESET_TOKEN` | 400 | Token reset password tidak valid/expired |
| `INVALID_CURRENT_PASSWORD` | 400 | Password saat ini salah |
| `INVALID_CREDENTIALS` | 401 | Email atau password salah |
| `TOKEN_MISSING` | 401 | Header Authorization tidak ditemukan |
| `INVALID_TOKEN` | 401 | Token JWT tidak valid |
| `TOKEN_EXPIRED` | 401 | Token JWT sudah kedaluwarsa |
| `USER_NOT_FOUND` | 404 | User tidak ditemukan |
| `DOCUMENT_NOT_FOUND` | 404 | Dokumen tidak ditemukan |
| `ADMIN_EMAIL_EXISTS` | 409 | Email admin sudah ada di whitelist |
| `USER_ALREADY_EXISTS` | 409 | Email user sudah terdaftar |
| `VALIDATION_FAILED` | 422 | Validasi input gagal |
| `DATABASE_ERROR` | 500 | Terjadi error operasi database |
| `INTERNAL_ERROR` | 500 | Terjadi error internal server |

---

## Standard Response Format

### Sukses
```json
{
  "success": true,
  "message": "Human-readable message",
  "data": { ... }
}
```

### Error
```json
{
  "success": false,
  "message": "Human-readable error message",
  "error": "ERROR_CODE",
  "details": "Human-readable error message: ERROR_CODE"
}
```
