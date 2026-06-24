# API Reference - DOMESv2

Dokumentasi lengkap API endpoints untuk DOMESv2 Backend.

## Base URL

| Environment | URL |
|-------------|-----|
| Development | `http://localhost:3000` |
| Production  | `https://domesv2.yourdomain.com` |

## Authentication

API menggunakan **JWT (JSON Web Token)** untuk autentikasi.
- Token didapat dari endpoint `/api/auth/login` atau `/api/auth/register`
- Kirim token di header: `Authorization: Bearer <token>`
- Token berlaku 24 jam (default)

---

## Ringkasan Endpoints

### Public Endpoints (Tanpa Auth)
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/` | Informasi dasar API |
| GET | `/api/health-check` | Cek status kesehatan sistem |
| POST | `/api/auth/register` | Pendaftaran user baru |
| POST | `/api/auth/login` | Login user |
| POST | `/api/auth/forgot-password` | Pengiriman email reset password |
| POST | `/api/auth/reset-password` | Reset password dengan token |
| GET | `/api/reference/:type` | Mengambil data referensi (agencies, sdgs, sectors, dll.) |
| GET | `/api/documents` | List dokumen terbitan publik |
| GET | `/api/documents/search` | Pencarian teks bebas dokumen |
| GET | `/api/documents/:id` | Detail dokumen berdasarkan ID atau slug |
| GET | `/api/documents/:id/related` | Rekomendasi dokumen terkait |
| GET | `/api/documents/:id/download` | Unduh file dokumen (dan tracking download) |
| GET | `/api/stats` | Statistik agregat platform |
| GET | `/api/analytics/:type` | Grafik analitik publik (overview, by-sdg, dll.) |
| POST | `/api/reports` | Mengirim laporan broken link |

### Protected Endpoints (Memerlukan JWT Bearer Token)
| Method | Endpoint | Role | Deskripsi |
|--------|----------|------|-----------|
| GET | `/api/user/me` | User & Admin | Profil user login |
| PUT | `/api/user/profile` | User & Admin | Edit profil diri |
| PUT | `/api/user/password` | User & Admin | Ganti password |
| GET | `/api/user/notifications` | User & Admin | Ambil preferensi notifikasi |
| PUT | `/api/user/notifications` | User & Admin | Edit preferensi notifikasi |
| GET | `/api/admin/emails` | Admin | List email whitelist admin |
| POST | `/api/admin/emails` | Admin | Tambah email ke whitelist admin |
| DELETE | `/api/admin/emails/:email` | Admin | Hapus email dari whitelist |
| GET | `/api/cms/dashboard` | User & Admin | Ringkasan dashboard editor |
| GET | `/api/cms/activity` | User & Admin | Aktivitas terbaru editor |
| GET | `/api/submissions` | User & Admin | List dokumen pengajuan editor |
| POST | `/api/submissions` | User & Admin | Submit final dokumen (Step 4) |
| POST | `/api/submissions/:id/draft` | User & Admin | Simpan draf dokumen (Step 1-3) |
| DELETE | `/api/submissions/:id` | User & Admin | Hapus dokumen pengajuan |
| PUT | `/api/submissions/:id/publish` | User & Admin | Terbitkan dokumen ke publik |
| PUT | `/api/submissions/:id/unpublish` | User & Admin | Tarik dokumen dari publik |
| GET | `/api/reports` | User & Admin | List laporan broken link |
| PUT | `/api/reports/:id/status` | User & Admin | Update status laporan link |
| GET | `/api/analytics/summary` | User & Admin | Ringkasan analitik internal |
| GET | `/api/analytics/top-downloads` | User & Admin | Top dokumen terunduh |
| GET | `/api/analytics/top-views` | User & Admin | Top dokumen dilihat |
| POST | `/api/upload` | User & Admin | Upload file dokumen/media |
| POST | `/api/upload/url-validate` | User & Admin | Validasi URL eksternal |
| POST | `/api/upload/avatar` | User & Admin | Upload foto profil/avatar |
| GET | `/api/users` | Admin | List user pengelola |
| POST | `/api/users` | Admin | Tambah user pengelola manual |
| PUT | `/api/users/:id` | Admin | Edit user pengelola |
| DELETE | `/api/users/:id` | Admin | Hapus user pengelola |

---

## Detail Endpoints

### 1. Authentication & Profiles

#### POST `/api/auth/register` (Public)
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

#### POST `/api/auth/login` (Public)
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

#### GET `/api/user/me` (Protected)
Mendapatkan data lengkap user yang sedang login beserta preferensinya.
* **Response 200:** Mengembalikan objek user lengkap.

#### PUT `/api/user/profile` (Protected)
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

#### PUT `/api/user/password` (Protected)
Mengganti password akun aktif.
* **Request Body:**
  ```json
  {
    "current_password": "rahasiaku123",
    "new_password": "passwordbaru123",
    "confirm_password": "passwordbaru123"
  }
  ```

#### GET & PUT `/api/user/notifications` (Protected)
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

#### GET `/api/admin/emails`
Mengambil daftar whitelist admin email.
* **Response 200:** Array objek email whitelist.

#### POST `/api/admin/emails`
Menambahkan email baru agar saat mendaftar otomatis menjadi Administrator.
* **Request Body:**
  ```json
  { "email": "new-admin@un.org" }
  ```

#### DELETE `/api/admin/emails/:email`
Menghapus email dari daftar whitelist.

---

### 3. Reference Data (Public)

#### GET `/api/reference/:type`
* Parameter `:type` yang valid: `agencies`, `sdgs`, `sectors`, `languages`, `joint-programmes`, `lnobs`, `non-un-partners`, `organizations`.
* **Response 200:** List objek data referensi (berisi id, code, name, icon/color jika ada).

---

### 4. Public Documents Discovery (Public)

#### GET `/api/documents`
Mendapatkan semua dokumen yang berstatus `published` dengan pagination & filter.
* **Query Parameters:**
  * `page` (default 1)
  * `limit` (default 10)
  * `agency` (code agensi)
  * `sdg` (code SDG)
  * `sector` (code sektor)
  * `language` (code bahasa)
  * `sort` (`newest`, `oldest`, `downloads`, `views`)

#### GET `/api/documents/search`
Pencarian teks bebas pada dokumen.
* **Query Parameters:** `q` (kata kunci pencarian), `sort` (`relevance`, `newest`, dsb.)

#### GET `/api/documents/:id`
Mencari dokumen berdasarkan ID numerik atau Slug teks unik.

#### GET `/api/documents/:id/related`
Mendapatkan rekomendasi dokumen lain yang memiliki irisan SDG atau sektor.

#### GET `/api/documents/:id/download`
Meningkatkan counter downloads dokumen dan mengembalikan link unduhan.

---

### 5. Broken Link Reporting (Public & Protected)

#### POST `/api/reports` (Public)
Mengajukan laporan link PDF dokumen yang rusak/404.
* **Request Body:**
  ```json
  {
    "document_id": 11,
    "reporter_name": "John Doe",
    "reporter_email": "johndoe@example.com",
    "details": "Tautan download PDF mengarah ke halaman kosong."
  }
  ```

#### GET `/api/reports` (Protected)
Mengambil daftar laporan yang diajukan. Query param: `status` (`all`, `pending`, `in_progress`, `resolved`).

#### PUT `/api/reports/:id/status` (Protected)
Memperbarui status penanganan laporan broken link.
* **Request Body:**
  ```json
  { "status": "resolved" } // Pilihan: pending, in_progress, resolved
  ```

---

### 6. CMS & Submissions Management (Protected)

#### GET `/api/cms/dashboard` & `/api/cms/activity`
Statistik performa CMS internal editor dan log aktivitas riwayat aksi terbaru.

#### POST `/api/submissions/:id/draft`
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

#### POST `/api/submissions`
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

#### PUT `/api/submissions/:id/publish` & `/api/submissions/:id/unpublish`
Mengubah visibilitas dokumen (memublikasikan ke portal publik atau menyembunyikan).

---

### 7. File Upload & Validation (Protected)

#### POST `/api/upload`
Mengunggah file. Menggunakan parser multipart form-data.
* **Request Form:**
  * `file`: File media (PDF, Word, JPG, PNG)
  * `type`: `document` (untuk PDF/Word) atau `cover` (untuk image)
* **Response 201:** `{"success": true, "url": "/uploads/random-uuid.pdf", "size": "1.2 MB"}`

#### POST `/api/upload/url-validate`
Validasi URL eksternal apakah merespons dengan HTTP status 200 OK.
* **Request Body:** `{ "url": "https://active-link.com/document.pdf" }`

#### POST `/api/upload/avatar`
Upload foto profil diri (avatar) user aktif. Form key: `avatar` (image).

---

### 8. CMS User Management (Protected - Admin Only)

#### GET, POST, PUT, DELETE pada `/api/users`
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
