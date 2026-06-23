# API Reference

## Base URL
Development: `http://localhost:3000`

## Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/` | No | Welcome message & info |
| GET | `/api/health-check` | No | System health check |
| POST | `/api/auth/register` | No | Register user baru |
| POST | `/api/auth/login` | No | Login dan dapatkan JWT |
| POST | `/api/auth/forgot-password` | No | Kirim email reset password |
| POST | `/api/auth/reset-password` | No | Reset password dengan token |
| GET | `/api/user/me` | Yes (JWT) | Profil user yang login |

---

## GET `/`

Welcome endpoint dengan informasi dasar API.

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

## GET `/api/health-check`

Pengecekan status database, redis, dan aplikasi.

**Response 200 (Healthy):**
```json
{
  "success": true,
  "message": "All systems operational",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T00:00:00Z",
    "services": {
      "application": "healthy",
      "database": "healthy",
      "redis": "disabled"
    }
  }
}
```

---

## POST `/api/auth/register`

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

**Response 201 (Success):**
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
      "created_at": "...",
      "updated_at": "..."
    }
  }
}
```

---

## POST `/api/auth/login`

Login dengan email dan password.

**Request Body:**
```json
{
  "email": "erlangga@un.org",
  "password": "password123",
  "captcha": "google-recaptcha-response-token"
}
```

**Response 200 (Success):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { "...user object..." }
  }
}
```

---

## POST `/api/auth/forgot-password`

Mengirim email reset password.

**Request Body:**
```json
{
  "email": "erlangga@un.org",
  "captcha": "google-recaptcha-response-token"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "If the email exists, a reset link has been sent",
  "data": null
}
```

> Endpoint selalu return 200 untuk mencegah email enumeration.

---

## POST `/api/auth/reset-password`

Mereset password dengan token dari email.

**Request Body:**
```json
{
  "token": "a1b2c3d4...64-hex-chars",
  "password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "Password has been reset successfully",
  "data": null
}
```

---

## GET `/api/user/me`

Profil user yang login.

**Headers:** `Authorization: Bearer <token>`

**Response 200:**
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

## Error Codes

| Kode | HTTP | Deskripsi |
|------|------|-----------|
| `INVALID_REQUEST_BODY` | 400 | Format body salah |
| `INVALID_RESET_TOKEN` | 400 | Token reset invalid/expired |
| `INVALID_CREDENTIALS` | 401 | Email/password salah |
| `TOKEN_MISSING` | 401 | Header Authorization tidak ada |
| `INVALID_TOKEN` | 401 | Token JWT tidak valid |
| `USER_NOT_FOUND` | 404 | User tidak ditemukan |
| `USER_ALREADY_EXISTS` | 409 | Email sudah terdaftar |
| `CAPTCHA_INVALID` | 422 | Captcha tidak valid |
| `CAPTCHA_MISSING` | 422 | Captcha tidak dikirim |
| `VALIDATION_FAILED` | 422 | Validasi input gagal |
| `DATABASE_ERROR` | 500 | Error database |
| `INTERNAL_ERROR` | 500 | Error internal |
