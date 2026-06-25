# API Contracts — DOMES V2

Dokumentasi lengkap endpoint API, method, request/response payload untuk backend DOMES V2.

---

## Base URL

| Environment | URL |
|-------------|-----|
| Development | `http://localhost:3000` |
| Production  | `https://domesv2.yourdomain.com` |

## Authentication

Seluruh endpoint (kecuali dinyatakan public) menggunakan **JWT Bearer Token** pada Header Authorization.

```
Authorization: Bearer <token>
```

## Standard Response Format

Setiap response API dibungkus oleh envelope standard.

### Success
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
  "details": "Detailed error string/validation logs"
}
```

### Paginated Response
```json
{
  "success": true,
  "message": "...",
  "data": {
    "items": [ ... ],
    "pagination": {
      "page": 1,
      "limit": 12,
      "totalItems": 1248,
      "totalPages": 104
    }
  }
}
```

---

## A. Authentication & Profiles (Public & Auth Required)

---

### POST /api/v2/auth/register (Public)

Mendaftarkan user baru ke sistem. Jika email yang didaftarkan terdaftar dalam whitelist email admin, user otomatis memiliki role `administrator`. Jika tidak, role default adalah `editor` / `viewer`.

**Request Body:**
```json
{
  "first_name": "Erlangga",
  "last_name": "Agustino",
  "position": "Administrator",
  "organization": "UNITED NATIONS",
  "phone_number": "+628123456789",
  "email": "erlangga@un.org",
  "password": "password123",
  "confirm_password": "password123",
  "captcha": "google-recaptcha-response-token"
}
```

| Field | Required | Validasi / Deskripsi |
|-------|----------|----------------------|
| `first_name` | ✅ | Nama depan user |
| `last_name` | ✅ | Nama belakang user |
| `email` | ✅ | Format email valid, harus unik |
| `password` | ✅ | Minimal 6 karakter |
| `confirm_password` | ✅ | Harus sama dengan `password` |
| `position` | ✅ | Jabatan user |
| `organization` | ✅ | Organisasi / instansi user |
| `phone_number` | ❌ | Nomor telepon |
| `captcha` | ❌* | Wajib jika diaktifkan di konfigurasi server (`RECAPTCHA_ENABLED=true`) |

**Response 201 (Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "erlangga",
      "name": "Erlangga Agustino",
      "first_name": "Erlangga",
      "last_name": "Agustino",
      "email": "erlangga@un.org",
      "type": "admin",
      "role": "administrator",
      "status": "active",
      "is_active": true,
      "position": "Administrator",
      "organization": "UNITED NATIONS",
      "phone_number": "+628123456789",
      "avatar_url": null,
      "registration_id": "f516df5d-4f1a-4d22-861c-843de9cc1e2e",
      "metadata": null,
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00"
    }
  }
}
```

**Response 409 (Conflict):**
```json
{
  "success": false,
  "message": "User with this email already exists",
  "error": "USER_ALREADY_EXISTS",
  "details": "User with this email already exists: USER_ALREADY_EXISTS"
}
```

**Response 422 (Validation Failed):**
```json
{
  "success": false,
  "message": "Passwords do not match",
  "error": "VALIDATION_FAILED",
  "details": "Passwords do not match: VALIDATION_FAILED"
}
```

---

### POST /api/v2/auth/login (Public)

Melakukan autentikasi menggunakan email dan password untuk mendapatkan JWT Token. Menolak autentikasi jika user telah dinonaktifkan (`is_active` bernilai `false`).

**Request Body:**
```json
{
  "email": "erlangga@un.org",
  "password": "password123",
  "captcha": "google-recaptcha-response-token"
}
```

| Field | Required | Validasi / Deskripsi |
|-------|----------|----------------------|
| `email` | ✅ | Format email valid |
| `password` | ✅ | Password akun |
| `captcha` | ❌* | Wajib jika diaktifkan di konfigurasi server |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "erlangga",
      "name": "Erlangga Agustino",
      "first_name": "Erlangga",
      "last_name": "Agustino",
      "email": "erlangga@un.org",
      "type": "admin",
      "role": "administrator",
      "status": "active",
      "is_active": true,
      "position": "Administrator",
      "organization": "UNITED NATIONS",
      "phone_number": "+628123456789",
      "avatar_url": "/uploads/avatars/erlangga.jpg",
      "registration_id": "f516df5d-4f1a-4d22-861c-843de9cc1e2e",
      "metadata": null,
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00"
    }
  }
}
```

**Response 401 (Unauthorized / Deactivated):**
```json
{
  "success": false,
  "message": "User account is deactivated",
  "error": "USER_DEACTIVATED",
  "details": "User account is deactivated: USER_DEACTIVATED"
}
```

---

### POST /api/v2/auth/forgot-password (Public)

Mengirimkan email instruksi reset password. Selalu mengembalikan HTTP 200 untuk mencegah enumerasi email.

**Request Body:**
```json
{
  "email": "erlangga@un.org",
  "captcha": "google-recaptcha-response-token"
}
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "If the email exists, a reset link has been sent",
  "data": null
}
```

---

### POST /api/v2/auth/reset-password (Public)

Melakukan reset password lama dengan password baru menggunakan reset token yang dikirim via email.

**Request Body:**
```json
{
  "token": "3a7b8e5c1d0f...",
  "password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Password has been reset successfully",
  "data": null
}
```

**Response 400 (Bad Request):**
```json
{
  "success": false,
  "message": "Invalid or expired reset token",
  "error": "INVALID_RESET_TOKEN",
  "details": "Invalid or expired reset token: INVALID_RESET_TOKEN"
}
```

---

### GET /api/v2/user/me (Auth Required)

Mengambil data profil lengkap dari user yang saat ini sedang login beserta preferensi notifikasinya.

**Headers:**
```
Authorization: Bearer <token>
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "User profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "erlangga",
    "name": "Erlangga Agustino",
    "first_name": "Erlangga",
    "last_name": "Agustino",
    "email": "erlangga@un.org",
    "type": "admin",
    "role": "administrator",
    "status": "active",
    "is_active": true,
    "position": "Administrator",
    "organization": "UNITED NATIONS",
    "phone_number": "+628123456789",
    "avatar_url": "/uploads/avatars/erlangga.jpg",
    "notification_preferences": {
      "id": "8fa57cbd-1334-4591-9fc7-2ef9da135011",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "document_approvals": true,
      "broken_link_reports": true,
      "system_updates": false,
      "email_notifications": true
    },
    "created_at": "2026-06-25T10:00:00+07:00",
    "updated_at": "2026-06-25T10:00:00+07:00"
  }
}
```

---

### PUT /api/v2/user/profile (Auth Required)

Mengubah data informasi diri dari user yang sedang login.

**Request Body:**
```json
{
  "first_name": "Erlangga",
  "last_name": "Agustino",
  "position": "Senior Administrator",
  "organization": "UNITED NATIONS (UNDP)",
  "phone_number": "+628123456789"
}
```

| Field | Required | Deskripsi |
|-------|----------|-----------|
| `first_name` | ✅ | Nama depan |
| `last_name` | ✅ | Nama belakang |
| `position` | ❌ | Jabatan baru |
| `organization` | ❌ | Organisasi baru |
| `phone_number` | ❌ | Telepon baru |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "first_name": "Erlangga",
    "last_name": "Agustino",
    "name": "Erlangga Agustino",
    "email": "erlangga@un.org",
    "position": "Senior Administrator",
    "organization": "UNITED NATIONS (UNDP)",
    "phone_number": "+628123456789",
    "is_active": true,
    "updated_at": "2026-06-25T11:00:00+07:00"
  }
}
```

---

### PUT /api/v2/user/password (Auth Required)

Mengubah password user.

**Request Body:**
```json
{
  "current_password": "password123",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Password changed successfully",
  "data": null
}
```

**Response 400 (Bad Request):**
```json
{
  "success": false,
  "message": "Current password is incorrect",
  "error": "INVALID_CURRENT_PASSWORD",
  "details": "Current password is incorrect: INVALID_CURRENT_PASSWORD"
}
```

---

### GET /api/v2/user/notifications (Auth Required)

Mengambil data preferensi notifikasi user.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Notification preferences retrieved successfully",
  "data": {
    "id": "8fa57cbd-1334-4591-9fc7-2ef9da135011",
    "created_at": "2026-06-25T10:00:00+07:00",
    "updated_at": "2026-06-25T10:00:00+07:00",
    "created_by": "System",
    "updated_by": "System",
    "is_active": true,
    "document_approvals": true,
    "broken_link_reports": true,
    "system_updates": false,
    "email_notifications": true
  }
}
```

---

### PUT /api/v2/user/notifications (Auth Required)

Mengubah preferensi notifikasi user.

**Request Body:**
```json
{
  "document_approvals": true,
  "broken_link_reports": false,
  "system_updates": true,
  "email_notifications": true
}
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Notification preferences updated successfully",
  "data": {
    "id": "8fa57cbd-1334-4591-9fc7-2ef9da135011",
    "created_at": "2026-06-25T10:00:00+07:00",
    "updated_at": "2026-06-25T11:30:00+07:00",
    "created_by": "System",
    "updated_by": "System",
    "is_active": true,
    "document_approvals": true,
    "broken_link_reports": false,
    "system_updates": true,
    "email_notifications": true
  }
}
```

---

## B. Public — Documents (Discovery)

Hanya dokumen yang memiliki `status` = `published` dan `is_active` = `true` yang akan tampil pada endpoint pencarian dan penemuan dokumen publik.

---

### GET /api/v2/documents (Public)

Mendapatkan daftar semua dokumen terbitan publik yang aktif dengan dukungan filtering, sorting, dan pagination.

**Query Parameters:**

| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `page` | integer | ❌ | Halaman aktif (Default: 1) |
| `limit` | integer | ❌ | Jumlah item per halaman (Default: 12, Max: 50) |
| `sort` | string | ❌ | `newest`, `oldest`, `downloads`, `views` (Default: `newest`) |
| `agencies` | string | ❌ | Koma terpisah kode agensi: `UNDP,UNICEF` |
| `sdgs` | string | ❌ | Koma terpisah kode SDG: `GOAL 1,GOAL 13` |
| `sectors` | string | ❌ | Koma terpisah kode sektor: `agriculture-food,economic-development` |
| `langs` | string | ❌ | Koma terpisah kode bahasa: `english,bahasa` |
| `yearFrom` | integer | ❌ | Batas awal tahun publikasi |
| `yearTo` | integer | ❌ | Batas akhir tahun publikasi |
| `jointProgrammes` | string | ❌ | Koma terpisah kode joint programme |
| `lnobs` | string | ❌ | Koma terpisah kode LNOB group: `women-girls` |
| `nonUnPartners` | string | ❌ | Koma terpisah kode non-un partner: `government` |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Documents retrieved successfully",
  "data": {
    "items": [
      {
        "id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
        "title": "Digital Economy and Financial Inclusion in Rural Indonesia",
        "slug": "digital-economy-financial-inclusion-rural-indonesia",
        "description": "This comprehensive report analyzes the rapid expansion of digital financial services across rural Indonesia...",
        "agency": "United Nations Development Programme",
        "year": 2026,
        "language": "English, Bahasa Indonesia",
        "file_size": "4.2 MB",
        "total_pages": 120,
        "type": "Report",
        "pub_status": "Published",
        "cover_image": "/uploads/covers/doc_001.jpg",
        "sdgs": ["GOAL 1", "GOAL 5", "GOAL 8", "GOAL 10"],
        "tags": ["digital economy", "financial inclusion", "fintech"],
        "views": 1234,
        "downloads": 567,
        "is_active": true,
        "created_at": "2026-06-25T10:00:00+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 12,
      "totalItems": 1,
      "totalPages": 1
    }
  }
}
```

---

### GET /api/v2/documents/search (Public)

Pencarian dokumen aktif dengan free-text query (termasuk fitur highlight kata kunci dan suggestions).

**Query Parameters:**

| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `q` | string | ✅ | Kata kunci pencarian |
| `page` | integer | ❌ | Halaman aktif (Default: 1) |
| `limit` | integer | ❌ | Jumlah item (Default: 12) |
| `sort` | string | ❌ | Pengurutan data (Default: `newest`) |
| `agencies`, `sdgs`, `sectors`, `langs`, `yearFrom`, `yearTo`, `jointProgrammes`, `lnobs`, `nonUnPartners` | string/int | ❌ | Filter opsional tambahan |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Search results retrieved successfully",
  "data": {
    "items": [
      {
        "id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
        "title": "Digital Economy and Financial Inclusion in Rural Indonesia",
        "slug": "digital-economy-financial-inclusion-rural-indonesia",
        "description": "This comprehensive report analyzes the rapid expansion of digital financial services across rural Indonesia...",
        "agency": "United Nations Development Programme",
        "year": 2026,
        "language": "English, Bahasa Indonesia",
        "file_size": "4.2 MB",
        "total_pages": 120,
        "type": "Report",
        "pub_status": "Published",
        "cover_image": "/uploads/covers/doc_001.jpg",
        "sdgs": ["GOAL 1", "GOAL 5", "GOAL 8", "GOAL 10"],
        "tags": ["digital economy", "financial inclusion", "fintech"],
        "views": 1234,
        "downloads": 567,
        "is_active": true,
        "highlight": {
          "title": "Digital Economy and <mark>Financial Inclusion</mark> in Rural Indonesia",
          "description": "...analyzes the rapid expansion of digital <mark>financial</mark> services..."
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 12,
      "totalItems": 1,
      "totalPages": 1
    },
    "suggestions": ["Green Economy", "Carbon Emission", "SDGs", "Paris Agreement"]
  }
}
```

---

### GET /api/v2/documents/{id} (Public)

Mengembalikan detail lengkap dari suatu dokumen publik yang aktif berdasarkan UUID v4 ID atau Slug uniknya.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Document retrieved successfully",
  "data": {
    "id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
    "code": "UNDP-2026-001",
    "slug": "digital-economy-financial-inclusion-rural-indonesia",
    "title": "Digital Economy and Financial Inclusion in Rural Indonesia",
    "agency": "UNDP",
    "year": 2026,
    "language": "English, Bahasa Indonesia",
    "file_url": "/uploads/documents/doc_001.pdf",
    "file_size": "4.2 MB",
    "date_added": "2026-06-25",
    "type": "Report",
    "total_pages": 120,
    "pub_status": "Published",
    "cover_image": "/uploads/covers/doc_001.jpg",
    "abstract": "This comprehensive report analyzes the rapid expansion of digital financial services across rural Indonesia. It highlights the profound impact of mobile banking and fintech solutions on local micro-economies.",
    "summary": "<b>Executive Overview</b><br><br>This extensive report provides an in-depth analysis...",
    "sdgs": [
      { "code": "GOAL 1", "name": "No Poverty", "icon": "/images/SDG-logos/SDG-1.png" },
      { "code": "GOAL 5", "name": "Gender Equality", "icon": "/images/SDG-logos/SDG-5.png" },
      { "code": "GOAL 8", "name": "Decent Work and Economic Growth", "icon": "/images/SDG-logos/SDG-8.png" },
      { "code": "GOAL 10", "name": "Reduced Inequalities", "icon": "/images/SDG-logos/SDG-10.png" }
    ],
    "tags": ["digital economy", "financial inclusion", "fintech", "women empowerment", "rural development"],
    "thematic_areas": ["Inclusive Economic Transformation"],
    "sectors": ["Economic Development", "Innovation and Technology", "Rural and Regional Development"],
    "lnob_groups": ["Women and Girls", "Rural populations"],
    "classification": {
      "lead_agency": "UNDP",
      "other_agencies": ["World Bank"],
      "joint_programme": "Climate Village Project (PROKLIM)",
      "geographic_scope": "National (Indonesia)",
      "non_un_partners": [
        { "type": "government", "name": "Ministry of Villages" }
      ]
    },
    "focal_point": {
      "name": "Budi Santoso",
      "email": "b.santoso@undp.org",
      "phone": "+62 812 3456 7890",
      "department": "Inclusive Growth Unit"
    },
    "views": 1234,
    "downloads": 567,
    "created_at": "2026-06-25T10:00:00+07:00",
    "updated_at": "2026-06-25T10:00:00+07:00",
    "is_active": true,
    "supporting_files": [
      { "url": "/uploads/supporting/appendix_a.pdf", "type": "appendix", "description": "Appendix A: Data Tables" }
    ]
  }
}
```

**Response 404 (Not Found):**
```json
{
  "success": false,
  "message": "Document not found",
  "error": "DOCUMENT_NOT_FOUND",
  "details": "Document not found: DOCUMENT_NOT_FOUND"
}
```

---

### GET /api/v2/documents/{id}/related (Public)

Mendapatkan rekomendasi dokumen lain yang memiliki kesamaan SDG atau Sektor dan bertatus aktif.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Related documents retrieved successfully",
  "data": [
    {
      "id": "b9c7df2a-29cc-4a92-b6ff-6fcf01ee1001",
      "code": "UNEP-2026-002",
      "slug": "climate-change-adaptation-coastal-communities",
      "title": "Climate Change Adaptation in Coastal Communities",
      "agency": "UNEP",
      "year": 2026,
      "cover_image": "/uploads/covers/doc_002.jpg",
      "is_active": true,
      "sdgs": [
        { "code": "GOAL 13", "name": "Climate Action", "icon": "/images/SDG-logos/SDG-13.png" }
      ]
    }
  ]
}
```

---

### GET /api/v2/documents/{id}/download (Public)

Men-track unduhan dokumen, menambah counter downloads, dan mengembalikan signed link/temporary download link file PDF bersangkutan.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Download link generated successfully",
  "data": {
    "download_url": "/uploads/documents/doc_001.pdf",
    "filename": "UNDP-2026-001_Digital_Economy_and_Financial_Inclusion_in_Rural_Indonesia.pdf",
    "file_size": "4.2 MB",
    "expires_at": "2026-06-25T11:00:00+07:00"
  }
}
```

---

## C. Public — Stats & Analytics

---

### GET /api/v2/stats (Public)

Mengembalikan statistik global platform (misalnya total dokumen, total agensi kontributor) untuk di-render pada portal Landing Page (Insights Banner).

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Platform stats retrieved successfully",
  "data": {
    "total_documents": 1796,
    "total_agencies": 24,
    "total_partners": 35,
    "total_sdg_goals": 17
  }
}
```

---

### GET /api/v2/analytics/overview (Public)

Mengembalikan metrik agregasi platform publik untuk dashboard analitik luar.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Analytics overview retrieved successfully",
  "data": {
    "total_documents": 12457,
    "active_agencies": 24,
    "monthly_downloads": 84200,
    "total_views": 456000,
    "total_downloads": 189000
  }
}
```

---

### GET /api/v2/analytics/uploads-over-time (Public)

Mengembalikan grafik garis total upload dokumen per tahun.

**Query Parameters:**
| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `fromYear` | integer | ❌ | Awal tahun (Default: 2019) |
| `toYear` | integer | ❌ | Akhir tahun (Default: 2024) |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Uploads over time analytics retrieved successfully",
  "data": [
    { "year": 2022, "count": 1890 },
    { "year": 2023, "count": 2340 },
    { "year": 2024, "count": 2670 },
    { "year": 2026, "count": 3120 }
  ]
}
```

---

### GET /api/v2/analytics/by-sdg (Public)

Mengembalikan total dokumen per SDG untuk kebutuhan visualisasi Chart batang.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Analytics by SDG retrieved successfully",
  "data": [
    { "sdg": "GOAL 1", "name": "No Poverty", "count": 1245, "color": "#E5243B" },
    { "sdg": "GOAL 13", "name": "Climate Action", "count": 1890, "color": "#3F7E44" }
  ]
}
```

---

### GET /api/v2/analytics/by-agency (Public)

Mengembalikan total dokumen terbitan per UN Agency.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Analytics by agency retrieved successfully",
  "data": [
    { "agency": "UNDP", "count": 1876 },
    { "agency": "UNEP", "count": 987 },
    { "agency": "UNICEF", "count": 1567 }
  ]
}
```

---

### GET /api/v2/analytics/by-sector (Public)

Mengembalikan total dokumen per sektor untuk Pie Chart.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Analytics by sector retrieved successfully",
  "data": [
    { "sector": "Environment and Climate Change", "count": 2345 },
    { "sector": "Economic Development", "count": 1234 }
  ]
}
```

---

### GET /api/v2/analytics/by-language (Public)

Mengembalikan pembagian total dokumen per bahasa untuk Donut Chart.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Analytics by language retrieved successfully",
  "data": [
    { "language": "English", "count": 8900 },
    { "language": "Bahasa Indonesia", "count": 5200 }
  ]
}
```

---

## D. Public — Broken Link Reports

---

### POST /api/v2/reports (Public)

Melaporkan link PDF yang rusak atau 404 dari dokumen publik.

**Request Body:**
```json
{
  "document_id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
  "reporter_name": "Budi Santoso",
  "reporter_email": "budi@example.com",
  "details": "The PDF link leads to a 404 error page.",
  "captcha": "google-recaptcha-response-token"
}
```

| Field | Required | Validasi / Deskripsi |
|-------|----------|----------------------|
| `document_id` | ✅ | ID Dokumen UUID v4 yang rusak |
| `reporter_name` | ✅ | Nama pelapor |
| `reporter_email` | ✅ | Email pelapor valid |
| `details` | ✅ | Keterangan kerusakan/deskripsi |
| `captcha` | ❌* | Diperlukan jika konfigurasi RECAPTCHA aktif |

**Response 201 (Created):**
```json
{
  "success": true,
  "message": "Report submitted successfully",
  "data": {
    "id": "9ca10cbd-1334-4591-9fc7-2ef9da135019",
    "document_id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
    "status": "open",
    "is_active": true,
    "created_at": "2026-06-25T10:00:00+07:00"
  }
}
```

---

## E. Master Data (Public)

Mengembalikan data statis dropdown isian form di Front End yang **aktif** (`is_active` = `true`). Setiap objek referensi memiliki struktur `V2Base` dengan UUID v4 ID.

---

### GET /api/v2/master/agencies (Public)

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Agencies retrieved successfully",
  "data": [
    {
      "id": "bfa86cfa-5b12-4cfb-bfe3-aa837df21601",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "UNDP",
      "name": "United Nations Development Programme",
      "logo_url": "/images/agency-logos/undp.png"
    }
  ]
}
```

---

### GET /api/v2/master/sdgs (Public)

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "SDGs retrieved successfully",
  "data": [
    {
      "id": "cfa86cfa-5b12-4cfb-bfe3-aa837df21602",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "GOAL 1",
      "name": "No Poverty",
      "icon": "/images/SDG-logos/SDG-1.png",
      "color": "#E5243B"
    }
  ]
}
```

---

### GET /api/v2/master/sectors (Public)

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Sectors retrieved successfully",
  "data": [
    {
      "id": "dfa86cfa-5b12-4cfb-bfe3-aa837df21603",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "economic-development",
      "name": "Economic Development"
    }
  ]
}
```

---

### GET /api/v2/master/languages (Public)

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Languages retrieved successfully",
  "data": [
    {
      "id": "efa86cfa-5b12-4cfb-bfe3-aa837df21604",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "english",
      "name": "English"
    }
  ]
}
```

---

### GET /api/v2/master/joint-programmes (Public)

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Joint programmes retrieved successfully",
  "data": [
    {
      "id": "ffa86cfa-5b12-4cfb-bfe3-aa837df21605",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "proklim",
      "name": "Climate Village Project (PROKLIM)"
    }
  ]
}
```

---

### GET /api/v2/master/lnobs (Public)

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "LNOB groups retrieved successfully",
  "data": [
    {
      "id": "0fa86cfa-5b12-4cfb-bfe3-aa837df21606",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "women-girls",
      "name": "Women and Girls"
    }
  ]
}
```

---

### GET /api/v2/master/non-un-partners (Public)

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Non-UN partner types retrieved successfully",
  "data": [
    {
      "id": "1fa86cfa-5b12-4cfb-bfe3-aa837df21607",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "government",
      "name": "Government"
    }
  ]
}
```

---

### GET /api/v2/master/organizations (Public)

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Organizations retrieved successfully",
  "data": [
    {
      "id": "2fa86cfa-5b12-4cfb-bfe3-aa837df21608",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "united-nations",
      "name": "UNITED NATIONS"
    }
  ]
}
```

---

### GET /api/v2/master/thematic-areas (Public)

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Thematic areas retrieved successfully",
  "data": [
    {
      "id": "3fa86cfa-5b12-4cfb-bfe3-aa837df21609",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "inclusive-economic-transformation",
      "name": "Inclusive Economic Transformation"
    }
  ]
}
```

---

## F. CMS — Dashboard (Auth Required)

---

### GET /api/v2/cms/dashboard (Auth Required)

Mengambil data statistik ringkas dashboard editor CMS.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Dashboard stats retrieved successfully",
  "data": {
    "total_documents": { "value": 1248, "change": 12.5, "trend": "up" },
    "total_views": { "value": 45200, "change": 8.3, "trend": "up" },
    "total_downloads": { "value": 8930, "change": -2.1, "trend": "down" },
    "total_users": { "value": 156, "change": 5.7, "trend": "up" },
    "pending_approvals": { "value": 23, "change": 0, "trend": "neutral" },
    "reports": { "value": 7, "change": -1, "trend": "down" }
  }
}
```

---

### GET /api/v2/cms/activity (Auth Required)

Mengambil daftar log aktivitas terbaru yang dilakukan para kontributor sistem.

**Query Parameters:**
| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `limit` | integer | ❌ | Jumlah aktivitas yang ingin ditarik (Default: 10) |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Recent activity retrieved successfully",
  "data": [
    {
      "id": 1,
      "type": "submission",
      "action": "created",
      "description": "New document submitted: 'Digital Economy Report 2024'",
      "user": "Erlangga Agustino",
      "user_avatar": "/uploads/avatars/erlangga.jpg",
      "timestamp": "2026-06-25T10:15:00+07:00",
      "time_ago": "2 minutes ago"
    }
  ]
}
```

---

## G. CMS — Submissions (Wizard & Mgmt - Auth Required)

---

### GET /api/v2/submissions (Auth Required)

Mengambil daftar berkas submissions (pengajuan) dokumen dengan status draf, pending, maupun terbit (menampilkan baik data yang aktif maupun dinonaktifkan).

**Query Parameters:**
| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `status` | string | ❌ | Filter status: `all`, `draft`, `pending_review`, `published`, `approved_unpublished` |
| `search` | string | ❌ | Cari berdasarkan judul dokumen |
| `page` | integer | ❌ | Default: 1 |
| `limit` | integer | ❌ | Default: 20 |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Submissions list retrieved successfully",
  "data": {
    "items": [
      {
        "id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
        "title": "Digital Economy and Financial Inclusion in Rural Indonesia",
        "short_description": "Analysis of digital financial services expansion...",
        "date_of_publication": "2026-06-15",
        "status": "published",
        "is_active": true,
        "agency": "United Nations Development Programme",
        "author": "Erlangga Agustino",
        "created_at": "2026-06-25T10:00:00+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "totalItems": 1,
      "totalPages": 1
    }
  }
}
```

---

### POST /api/v2/submissions (Auth Required)

Membuat pengajuan dokumen secara final (pada Step 4 Wizard). Bisa melampirkan parameter `is_active` jika ingin diubah status keaktifannya.

**Request Body (Lengkap):**
```json
{
  "title": "Digital Economy and Financial Inclusion in Rural Indonesia",
  "short_description": "Analysis of digital financial services expansion...",
  "abstract": "This comprehensive report analyzes...",
  "detailed_summary": "<b>Executive Overview</b><br><br>This extensive report...",
  "date_of_publication": "2026-06-15",
  "total_pages": 120,
  "language": "English, Bahasa Indonesia",
  "publication_status": "Published",
  "tags": ["digital economy", "financial inclusion", "fintech"],
  "file_url": "/uploads/documents/doc_001.pdf",
  "file_size": "4.2 MB",
  "cover_image_url": "/uploads/covers/doc_001.jpg",
  "external_url": "",
  "supporting_files": [
    { "url": "/uploads/supporting/appendix_a.pdf", "type": "appendix", "description": "Appendix A: Data Tables" }
  ],
  "agency": "UNDP",
  "focal_point": {
    "name": "Budi Santoso",
    "email": "b.santoso@undp.org",
    "phone": "+62 812 3456 7890",
    "department": "Inclusive Growth Unit"
  },
  "sdgs": ["GOAL 1", "GOAL 5", "GOAL 8", "GOAL 10"],
  "sectors": ["economic-development"],
  "lnob_groups": ["women-girls"],
  "joint_programme": "proklim",
  "other_agencies": ["World Bank"],
  "non_un_partners": [
    { "type": "government", "name": "Ministry of Villages" }
  ],
  "thematic_areas": ["Inclusive Economic Transformation"],
  "geographic_scope": "National (Indonesia)",
  "is_active": true
}
```

**Response 201 (Created):**
```json
{
  "success": true,
  "message": "Submission created successfully",
  "data": {
    "id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
    "code": "UNDP-2026-001",
    "slug": "digital-economy-financial-inclusion-rural-indonesia",
    "title": "Digital Economy and Financial Inclusion in Rural Indonesia",
    "status": "pending_review",
    "is_active": true,
    "created_at": "2026-06-25T10:00:00+07:00"
  }
}
```

---

### POST /api/v2/submissions/{id}/draft (Auth Required)

Menyimpan langkah draf (langkah 1-3) ke database untuk disimpan sementara. Bisa menambahkan `"is_active": false` untuk menonaktifkan sementara dari sistem.

**Request Body:**
```json
{
  "step": 2,
  "data": {
    "title": "Draf Dokumen Baru",
    "short_summary": "Summary draf...",
    "focal_point_name": "Budi Santoso",
    "focal_point_email": "b.santoso@undp.org",
    "is_active": true
  }
}
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Draft saved successfully",
  "data": {
    "id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
    "step": 2,
    "saved_at": "2026-06-25T10:20:00+07:00"
  }
}
```

---

### DELETE /api/v2/submissions/{id} (Auth Required)

Menghapus pengajuan dokumen berdasarkan UUID v4 ID.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Submission deleted successfully",
  "data": null
}
```

---

### PUT /api/v2/submissions/{id}/publish (Auth Required)

Mempublikasikan dokumen agar terbit dan terlihat di portal publik.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Document published successfully",
  "data": {
    "id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
    "status": "published",
    "published_at": "2026-06-25T10:30:00+07:00"
  }
}
```

---

### PUT /api/v2/submissions/{id}/unpublish (Auth Required)

Menarik kembali publikasi dokumen agar kembali ke status unpublished (CMS only).

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Document unpublished successfully",
  "data": {
    "id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
    "status": "approved_unpublished"
  }
}
```

---

## H. CMS — Master Management (Auth Required - Admin / Editor)

Menyediakan fungsi pengelolaan data pilihan master secara modular oleh Administrator. Mendukung set `is_active` = `false` agar data master disembunyikan dari dropdown publik namun datanya tidak hilang untuk relasi tabel yang sudah ada.

---

### GET /api/v2/cms/master/{type} (Auth Required)

Mengambil seluruh data master dari tipe tertentu (termasuk yang aktif maupun tidak).
* **Tipe valid (`{type}`):** `agencies`, `sdgs`, `sectors`, `languages`, `joint-programmes`, `lnobs`, `non-un-partners`, `organizations`, `thematic-areas`

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Master list retrieved successfully",
  "data": [
    {
      "id": "dfa86cfa-5b12-4cfb-bfe3-aa837df21603",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "code": "economic-development",
      "name": "Economic Development"
    }
  ]
}
```

---

### POST /api/v2/cms/master/{type} (Auth Required - Admin Only)

Menambahkan data master baru.

**Request Body:**
```json
{
  "code": "test-sector-xyz",
  "name": "Test Sector XYZ",
  "logo_url": "",
  "icon": "",
  "color": "",
  "is_active": true
}
```

**Response 201 (Created):**
```json
{
  "success": true,
  "message": "Master item created successfully",
  "data": {
    "id": "8fa86cfa-5b12-4cfb-bfe3-aa837df21612",
    "created_at": "2026-06-25T10:35:00+07:00",
    "updated_at": "2026-06-25T10:35:00+07:00",
    "created_by": "admin@un.org",
    "updated_by": "admin@un.org",
    "is_active": true,
    "code": "test-sector-xyz",
    "name": "Test Sector XYZ"
  }
}
```

---

### PUT /api/v2/cms/master/{type}/{code} (Auth Required - Admin Only)

Memperbarui data master berdasarkan parameter kode unik. Mendukung update parameter `is_active`.

**Request Body:**
```json
{
  "name": "Test Sector XYZ Updated",
  "logo_url": "/images/logo-updated.png",
  "is_active": false
}
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Master item updated successfully",
  "data": {
    "id": "8fa86cfa-5b12-4cfb-bfe3-aa837df21612",
    "created_at": "2026-06-25T10:35:00+07:00",
    "updated_at": "2026-06-25T10:40:00+07:00",
    "created_by": "admin@un.org",
    "updated_by": "admin@un.org",
    "is_active": false,
    "code": "test-sector-xyz",
    "name": "Test Sector XYZ Updated",
    "logo_url": "/images/logo-updated.png"
  }
}
```

---

### DELETE /api/v2/cms/master/{type}/{code} (Auth Required - Admin Only)

Menghapus item master secara permanen dari database.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Master item deleted successfully",
  "data": null
}
```

---

## I. CMS — Users Management (Auth Required - Admin Only)

---

### GET /api/v2/users (Auth Required - Admin Only)

Mengembalikan daftar user pengelola sistem.

**Query Parameters:**
| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `search` | string | ❌ | Cari berdasarkan nama, email, organisasi |
| `role` | string | ❌ | Filter role: `administrator`, `editor`, `viewer` |
| `status` | string | ❌ | Filter status: `active`, `inactive` |
| `page` | integer | ❌ | Halaman aktif (Default: 1) |
| `limit` | integer | ❌ | Jumlah item per halaman (Default: 20) |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": {
    "items": [
      {
        "id": 1,
        "first_name": "Erlangga",
        "last_name": "Agustino",
        "email": "erlangga@un.org",
        "phone_number": "+628123456789",
        "organization": "UNITED NATIONS",
        "position": "Administrator",
        "role": "administrator",
        "status": "active",
        "is_active": true,
        "avatar_url": "/uploads/avatars/erlangga.jpg",
        "created_at": "2026-06-25T10:00:00+07:00",
        "last_login": "2026-06-25T10:05:00+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "totalItems": 1,
      "totalPages": 1
    }
  }
}
```

---

### POST /api/v2/users (Auth Required - Admin Only)

Membuat user baru secara manual dari CMS panel.

**Request Body:**
```json
{
  "first_name": "Budi",
  "last_name": "Santoso",
  "email": "budi.santoso@un.org",
  "password": "password123",
  "confirm_password": "password123",
  "organization": "WHO",
  "position": "Health Officer",
  "phone_number": "+6281122334455",
  "role": "editor",
  "status": "active",
  "is_active": true
}
```

| Field | Required | Validasi / Deskripsi |
|-------|----------|----------------------|
| `first_name` | ✅ | Nama depan |
| `last_name` | ✅ | Nama belakang |
| `email` | ✅ | Email valid, belum terdaftar |
| `password` | ✅ | Minimal 6 karakter |
| `confirm_password` | ✅ | Harus sama dengan `password` |
| `organization` | ❌ | Nama organisasi |
| `position` | ❌ | Jabatan user |
| `phone_number` | ❌ | Nomor telepon |
| `role` | ✅ | Pilihan: `administrator`, `editor`, `viewer` |
| `status` | ❌ | Status aktif: `active`, `inactive` |
| `is_active` | ❌ | Flag aktif akun (Default: `true`) |

**Response 201 (Created):**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 2,
    "first_name": "Budi",
    "last_name": "Santoso",
    "email": "budi.santoso@un.org",
    "organization": "WHO",
    "position": "Health Officer",
    "role": "editor",
    "status": "active",
    "is_active": true,
    "created_at": "2026-06-25T10:45:00+07:00"
  }
}
```

---

### PUT /api/v2/users/{id} (Auth Required - Admin Only)

Memperbarui data pengelola sistem. Mendukung pembaruan parameter `is_active` (misal untuk deaktivasi akun editor).

**Request Body:**
```json
{
  "first_name": "Budi",
  "last_name": "Santoso",
  "organization": "WHO",
  "position": "Senior Health Officer",
  "phone_number": "+6281122334455",
  "role": "administrator",
  "status": "active",
  "is_active": false
}
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 2,
    "first_name": "Budi",
    "last_name": "Santoso",
    "email": "budi.santoso@un.org",
    "organization": "WHO",
    "position": "Senior Health Officer",
    "role": "administrator",
    "status": "active",
    "is_active": false,
    "updated_at": "2026-06-25T10:50:00+07:00"
  }
}
```

---

### DELETE /api/v2/users/{id} (Auth Required - Admin Only)

Menghapus user pengelola sistem.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": null
}
```

---

## J. CMS — Analytics (Auth Required)

---

### GET /api/v2/analytics/summary (Auth Required)

Mengambil ringkasan metrik analitik internal CMS (views, downloads, active users) dengan perioda waktu tertentu.

**Query Parameters:**
| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `period` | string | ❌ | Periode: `7d`, `30d`, `90d`, `1y` (Default: `30d`) |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Analytics summary retrieved successfully",
  "data": {
    "total_downloads": { "value": 24592, "change": 12.5, "trend": "up" },
    "total_views": { "value": 89401, "change": 8.2, "trend": "up" },
    "active_users": { "value": 3240, "change": -2.1, "trend": "down" }
  }
}
```

---

### GET /api/v2/analytics/top-downloads (Auth Required)

Mengambil daftar dokumen dengan total unduhan terbanyak beserta persentase progress visualnya.

**Query Parameters:**
| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `limit` | integer | ❌ | Jumlah dokumen teratas (Default: 10) |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Top downloads retrieved successfully",
  "data": [
    { "title": "Digital Economy Report 2026", "downloads": 1234, "progress": 100 },
    { "title": "Climate Change Adaptation Guide", "downloads": 987, "progress": 80 }
  ]
}
```

---

### GET /api/v2/analytics/top-views (Auth Required)

Mengambil daftar dokumen yang paling sering dilihat/dikunjungi.

**Query Parameters:**
| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `limit` | integer | ❌ | Jumlah dokumen teratas (Default: 10) |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Top views retrieved successfully",
  "data": [
    { "title": "Digital Economy Report 2026", "category": "UNDP", "views": 5678 },
    { "title": "Climate Change Adaptation Guide", "category": "UNEP", "views": 4567 }
  ]
}
```

---

## K. CMS — Reports (Auth Required)

---

### GET /api/v2/reports (Auth Required)

Mengambil seluruh daftar laporan broken links yang diajukan oleh user publik.

**Query Parameters:**
| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `status` | string | ❌ | Filter status: `all`, `open`, `in_progress`, `resolved` |
| `search` | string | ❌ | Cari berdasarkan kata kunci judul dokumen |
| `page` | integer | ❌ | Default: 1 |
| `limit` | integer | ❌ | Default: 20 |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Reports retrieved successfully",
  "data": {
    "items": [
      {
        "id": "9ca10cbd-1334-4591-9fc7-2ef9da135019",
        "document_id": "7da60cbd-1334-4591-9fc7-2ef9da135014",
        "document_title": "Digital Economy and Financial Inclusion in Rural Indonesia",
        "reporter_name": "Budi Santoso",
        "reporter_email": "budi@example.com",
        "details": "The PDF link leads to a 404 error page.",
        "status": "open",
        "is_active": true,
        "created_at": "2026-06-25T10:00:00+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "totalItems": 1,
      "totalPages": 1
    }
  }
}
```

---

### PUT /api/v2/reports/{id}/status (Auth Required)

Mengubah status penyelesaian laporan kerusakan link (e.g. dari `open` menjadi `in_progress` atau `resolved`).

**Request Body:**
```json
{
  "status": "in_progress"
}
```

| Field | Required | Validasi / Deskripsi |
|-------|----------|----------------------|
| `status` | ✅ | Pilihan status: `open`, `in_progress`, `resolved` |

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Report status updated successfully",
  "data": {
    "id": "9ca10cbd-1334-4591-9fc7-2ef9da135019",
    "status": "in_progress",
    "is_active": true,
    "updated_at": "2026-06-25T11:00:00+07:00"
  }
}
```

**Response 400 (Bad Request):**
```json
{
  "success": false,
  "message": "Invalid status value",
  "error": "VALIDATION_FAILED",
  "details": "Status must be one of: open, in_progress, resolved"
}
```

---

## L. CMS — Whitelist Settings (Auth Required - Admin Only)

---

### GET /api/v2/admin/emails (Auth Required - Admin Only)

Mengembalikan seluruh daftar alamat email whitelist admin. User yang melakukan registrasi dengan salah satu email dari list ini otomatis diposisikan sebagai administrator.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Admin emails retrieved successfully",
  "data": [
    {
      "id": "cca86cfa-5b12-4cfb-bfe3-aa837df21699",
      "created_at": "2026-06-25T10:00:00+07:00",
      "updated_at": "2026-06-25T10:00:00+07:00",
      "created_by": "System",
      "updated_by": "System",
      "is_active": true,
      "email": "admin@un.org",
      "added_at": "2026-06-25T10:00:00+07:00"
    }
  ]
}
```

---

### POST /api/v2/admin/emails (Auth Required - Admin Only)

Menambahkan alamat email baru ke dalam daftar whitelist admin.

**Request Body:**
```json
{
  "email": "newadmin@un.org"
}
```

| Field | Required | Validasi / Deskripsi |
|-------|----------|----------------------|
| `email` | ✅ | Format email valid, harus unik |

**Response 201 (Created):**
```json
{
  "success": true,
  "message": "Admin email added successfully",
  "data": {
    "id": "cca86cfa-5b12-4cfb-bfe3-aa837df21610",
    "created_at": "2026-06-25T11:15:00+07:00",
    "updated_at": "2026-06-25T11:15:00+07:00",
    "created_by": "admin@un.org",
    "updated_by": "admin@un.org",
    "is_active": true,
    "email": "newadmin@un.org",
    "added_at": "2026-06-25T11:15:00+07:00"
  }
}
```

---

### DELETE /api/v2/admin/emails/{email} (Auth Required - Admin Only)

Menghapus alamat email dari daftar whitelist admin.

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Admin email removed successfully",
  "data": null
}
```

---

## M. File Uploads & Validations (Auth Required)

---

### POST /api/v2/upload (Auth Required)

Mengunggah berkas dokumen utama, cover, pendukung, atau avatar pengelola. Menggunakan payload tipe `multipart/form-data`.

**Request:** `multipart/form-data`

| Field | Tipe | Required | Deskripsi |
|-------|------|----------|-----------|
| `file` | file | ✅ | Berkas file biner (Max 50MB) |
| `type` | string | ✅ | Tipe berkas: `document`, `cover`, `supporting`, `avatar` |
| `submission_id` | string | ❌ | ID Submission opsional (UUID v4) |

**Response 201 (Created):**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "url": "/uploads/3cfa86cf-5b12-4cfb-bfe3-aa837df21633.pdf",
    "filename": "3cfa86cf-5b12-4cfb-bfe3-aa837df21633.pdf",
    "original_name": "Digital_Economy_Report_2026.pdf",
    "file_size": "4.2 MB",
    "mime_type": "application/pdf"
  }
}
```

**Response 413 (Payload Too Large):**
```json
{
  "success": false,
  "message": "File size exceeds the maximum limit of 50MB",
  "error": "FILE_TOO_LARGE",
  "details": "File size exceeds the maximum limit of 50MB: FILE_TOO_LARGE"
}
```

---

### POST /api/v2/upload/url-validate (Auth Required)

Melakukan verifikasi eksternal link/URL dokumen apakah dapat diakses (berstatus HTTP 200 OK) atau error/terblokir.

**Request Body:**
```json
{
  "url": "https://example.com/documents/report.pdf"
}
```

**Response 200 (Valid):**
```json
{
  "success": true,
  "message": "URL is valid",
  "data": {
    "url": "https://example.com/documents/report.pdf",
    "accessible": true,
    "content_type": "application/pdf",
    "file_size": "3.5 MB"
  }
}
```

**Response 200 (Invalid / Not Accessible):**
```json
{
  "success": true,
  "message": "URL is not accessible",
  "data": {
    "url": "https://example.com/documents/report.pdf",
    "accessible": false,
    "error": "HTTP 404 Not Found"
  }
}
```

---

### POST /api/v2/upload/avatar (Auth Required)

Mengunggah foto profil / avatar baru untuk user yang sedang aktif.

**Request:** `multipart/form-data`

| Field | Tipe | Required | Deskripsi |
|-------|------|----------|-----------|
| `avatar` | file | ✅ | Berkas gambar (Max 2MB, format: jpg/png/webp) |

**Response 201 (Created):**
```json
{
  "success": true,
  "message": "Avatar uploaded successfully",
  "data": {
    "avatar_url": "/uploads/avatars/4cfa86cf-5b12-4cfb-bfe3-aa837df21634.jpg"
  }
}
```

---

## N. System (Public)

---

### GET /api/v2/health-check (Public)

Pengecekan status operasional dari database, koneksi Redis, dan status runtime backend aplikasi.

**Response 200 (Healthy):**
```json
{
  "success": true,
  "message": "All systems operational",
  "data": {
    "status": "healthy",
    "timestamp": "2026-06-25T10:00:00+07:00",
    "services": {
      "application": "healthy",
      "database": "healthy",
      "redis": "disabled"
    }
  }
}
```

**Response 503 (Service Unavailable):**
```json
{
  "success": true,
  "message": "Service degraded",
  "data": {
    "status": "unhealthy",
    "timestamp": "2026-06-25T10:00:00+07:00",
    "services": {
      "application": "healthy",
      "database": "unhealthy",
      "database_error": "database connection timed out",
      "redis": "disabled"
    }
  }
}
```

---

## Error Codes

Daftar parameter standard kesalahan internal API yang dikirimkan di dalam property payload `"error"`.

| Kode | HTTP Status | Deskripsi |
|------|-------------|-----------|
| `INVALID_REQUEST_BODY` | 400 | Format request body JSON salah atau tidak dapat di-parse |
| `INVALID_RESET_TOKEN` | 400 | Token reset password tidak valid atau telah expired |
| `INVALID_CURRENT_PASSWORD` | 400 | Password saat ini salah / tidak cocok |
| `INVALID_CREDENTIALS` | 401 | Email atau password tidak sesuai |
| `USER_DEACTIVATED` | 401 | Akun user telah dinonaktifkan (`is_active` = `false`) |
| `TOKEN_MISSING` | 401 | Header request `Authorization` JWT tidak ditemukan |
| `INVALID_TOKEN` | 401 | Tanda tangan Token JWT tidak valid / telah rusak |
| `TOKEN_EXPIRED` | 401 | Token JWT telah melewati batas waktu berlaku |
| `FORBIDDEN` | 403 | User tidak memiliki hak akses/role yang memadai |
| `USER_NOT_FOUND` | 404 | User tidak ditemukan di database |
| `DOCUMENT_NOT_FOUND` | 404 | Dokumen tidak ditemukan di database |
| `SUBMISSION_NOT_FOUND` | 404 | Submissions draf/pending tidak ditemukan |
| `ADMIN_EMAIL_NOT_FOUND` | 404 | Email admin tidak ditemukan dalam daftar whitelist |
| `USER_ALREADY_EXISTS` | 409 | Alamat email pendaftaran sudah digunakan |
| `ADMIN_EMAIL_EXISTS` | 409 | Alamat email whitelist admin sudah didaftarkan sebelumnya |
| `FILE_TOO_LARGE` | 413 | Ukuran file unggahan melebihi batas maximum (50MB / 2MB) |
| `CAPTCHA_INVALID` | 422 | Google Recaptcha verification token gagal divalidasi |
| `CAPTCHA_MISSING` | 422 | Verification token Recaptcha kosong / tidak dikirim |
| `VALIDATION_FAILED` | 422 | Pengisian data request tidak lolos kriteria validasi input |
| `DATABASE_ERROR` | 500 | Terjadi kegagalan operasi internal database engine |
| `INTERNAL_ERROR` | 500 | Terjadi kegagalan runtime internal server error |
