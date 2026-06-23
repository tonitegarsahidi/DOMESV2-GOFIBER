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

## Endpoints

| Method | Endpoint | Auth | Deskripsi |
|--------|----------|------|-----------|
| GET | `/` | ❌ | Informasi API |
| GET | `/api/health-check` | ❌ | Cek status sistem |
| POST | `/api/auth/register` | ❌ | Registrasi user baru |
| POST | `/api/auth/login` | ❌ | Login user |
| POST | `/api/auth/forgot-password` | ❌ | Lupa password |
| POST | `/api/auth/reset-password` | ❌ | Reset password |
| GET | `/api/user/me` | ✅ | Profil user login |

---

### GET `/`

Informasi dasar API.

**Response 200:**
```json
{
  "success": true,
  "message": "Hello Domes v2",
  "version": "1.0.0",
  "docs": "/api/health-check"
}
```

---

### GET `/api/health-check`

Pengecekan status database, redis, dan aplikasi.

**Response 200 (Sehat):**
```json
{
  "success": true,
  "message": "All systems operational",
  "data": {
    "status": "healthy",
    "timestamp": "2024-06-23T10:00:00+07:00",
    "services": {
      "application": "healthy",
      "database": "healthy",
      "redis": "disabled"
    }
  }
}
```

**Response 503 (Tidak Sehat):**
```json
{
  "status": "unhealthy",
  "timestamp": "2024-06-23T10:00:00+07:00",
  "services": {
    "application": "healthy",
    "database": "unhealthy",
    "database_error": "database not initialized",
    "redis": "disabled"
  }
}
```

---

### POST `/api/auth/register`

Mendaftarkan user baru.

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

**Validation Rules:**
| Field | Rule |
|-------|------|
| `first_name` | Required |
| `last_name` | Required |
| `email` | Required, format email valid |
| `password` | Required, minimal 6 karakter |
| `confirm_password` | Required, harus sama dengan `password` |
| `captcha` | Required jika reCAPTCHA diaktifkan |

**Response 201 (Sukses):**
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
      "type": null,
      "position": "Administrator",
      "organization": "UNITED NATIONS",
      "phone_number": "+628123456789",
      "created_at": "2024-06-23T10:00:00+07:00",
      "updated_at": "2024-06-23T10:00:00+07:00"
    }
  }
}
```

**Response 409 (Email sudah terdaftar):**
```json
{
  "success": false,
  "message": "User with this email already exists",
  "error": "USER_ALREADY_EXISTS",
  "details": "User with this email already exists: USER_ALREADY_EXISTS"
}
```

**Response 422 (Validasi gagal):**
```json
{
  "success": false,
  "message": "Passwords do not match",
  "error": "VALIDATION_FAILED",
  "details": "Passwords do not match: VALIDATION_FAILED"
}
```

---

### POST `/api/auth/login`

Login dengan email dan password.

**Request Body:**
```json
{
  "email": "erlangga@un.org",
  "password": "password123",
  "captcha": "google-recaptcha-response-token"
}
```

**Response 200 (Sukses):**
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
      "type": null,
      "position": "Administrator",
      "organization": "UNITED NATIONS",
      "phone_number": "+628123456789",
      "created_at": "2024-06-23T10:00:00+07:00",
      "updated_at": "2024-06-23T10:00:00+07:00"
    }
  }
}
```

**Response 401 (Kredensial salah):**
```json
{
  "success": false,
  "message": "Invalid credentials",
  "error": "INVALID_CREDENTIALS",
  "details": "Invalid credentials: INVALID_CREDENTIALS"
}
```

---

### POST `/api/auth/forgot-password`

Mengirim email reset password.

**Request Body:**
```json
{
  "email": "erlangga@un.org",
  "captcha": "google-recaptcha-response-token"
}
```

**Response 200 (Sukses):**
```json
{
  "success": true,
  "message": "If the email exists, a reset link has been sent",
  "data": null
}
```

> **Catatan:** Endpoint selalu return 200 meskipun email tidak ditemukan (keamanan - mencegah email enumeration).

---

### POST `/api/auth/reset-password`

Mereset password menggunakan token dari email.

**Request Body:**
```json
{
  "token": "a1b2c3d4e5f6...64-hex-char-token",
  "password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Validation Rules:**
| Field | Rule |
|-------|------|
| `token` | Required, token dari email reset password |
| `password` | Required, minimal 6 karakter |
| `confirm_password` | Required, harus sama dengan `password` |

**Response 200 (Sukses):**
```json
{
  "success": true,
  "message": "Password has been reset successfully",
  "data": null
}
```

**Response 400 (Token invalid/expired):**
```json
{
  "success": false,
  "message": "Invalid or expired reset token",
  "error": "INVALID_RESET_TOKEN",
  "details": "Invalid or expired reset token: INVALID_RESET_TOKEN"
}
```

---

### GET `/api/user/me`

Mendapatkan profil user yang sedang login.

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 200 (Sukses):**
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
    "position": "Administrator",
    "organization": "UNITED NATIONS",
    "phone_number": "+628123456789"
  }
}
```

**Response 401 (Token tidak valid):**
```json
{
  "success": false,
  "message": "Missing authorization header",
  "error": "TOKEN_MISSING",
  "details": "Missing authorization header: TOKEN_MISSING"
}
```

---

## Google reCAPTCHA v2 Integration

Beberapa endpoint (`/register`, `/login`, `/forgot-password`) dilindungi dengan Google reCAPTCHA v2.

### Setup di Google Cloud Console
1. Buka [Google Cloud Console](https://console.cloud.google.com/security/recaptcha)
2. Buat reCAPTCHA v2 → "I'm not a robot" checkbox
3. Dapatkan **Site Key** dan **Secret Key**

### Konfigurasi Environment
```env
# .env
RECAPTCHA_SITE_KEY=6Lc..._site_key_untuk_frontend
RECAPTCHA_SECRET_KEY=6Lc..._secret_key_untuk_backend
RECAPTCHA_ENABLED=true
```

### Pemasangan di Frontend
```html
<!-- Load reCAPTCHA JS -->
<script src="https://www.google.com/recaptcha/api.js" async defer></script>

<!-- Tempatkan widget di form login/register -->
<form>
  <input type="email" name="email" />
  <input type="password" name="password" />
  <div class="g-recaptcha" data-sitekey="{{ RECAPTCHA_SITE_KEY }}"></div>
  <button type="submit">Submit</button>
</form>
```

### Cara Kirim Token
Setelah user centang "I'm not a robot", Google memberikan token. Kirim token tersebut ke backend di field `captcha`:

```javascript
// Contoh dengan fetch
const token = grecaptcha.getResponse();

fetch('/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'secret',
    captcha: token   // <--- token dari reCAPTCHA
  })
});
```

### Nonaktifkan Captcha (Development)
Set `RECAPTCHA_ENABLED=false` di `.env` agar tidak perlu captcha saat development.

---

## Error Codes

| Kode | HTTP Status | Deskripsi |
|------|-------------|-----------|
| `INVALID_REQUEST_BODY` | 400 | Format request body salah |
| `INVALID_RESET_TOKEN` | 400 | Token reset password tidak valid/expired |
| `INVALID_CREDENTIALS` | 401 | Email atau password salah |
| `TOKEN_MISSING` | 401 | Header Authorization tidak ada |
| `INVALID_TOKEN` | 401 | Token JWT tidak valid |
| `TOKEN_EXPIRED` | 401 | Token JWT sudah expired |
| `USER_NOT_FOUND` | 404 | User tidak ditemukan |
| `USER_ALREADY_EXISTS` | 409 | Email sudah terdaftar |
| `CAPTCHA_INVALID` | 422 | Captcha tidak valid |
| `CAPTCHA_MISSING` | 422 | Captcha tidak dikirim |
| `VALIDATION_FAILED` | 422 | Validasi input gagal |
| `DATABASE_ERROR` | 500 | Error database |
| `INTERNAL_ERROR` | 500 | Error internal server |

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

### Pagination (jika ada)
Belum diimplementasikan.
