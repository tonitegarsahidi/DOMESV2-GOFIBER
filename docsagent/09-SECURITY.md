# Security Guide

## 1. JWT (JSON Web Token)

### Implementation
- **File:** `internal/middleware/jwt_middleware.go`
- **Algoritma:** HS256 (HMAC-SHA256)
- **Secret:** Dikonfigurasi via `JWT_SECRET` di `.env`
- **Expiry:** Default 24 jam (via `JWT_EXPIRES_IN`)

### JWT Claims
```go
claims := jwt.MapClaims{
    "user_id": user.ID,
    "email":   user.Email,
    "exp":     time.Now().Add(s.cfg.JWT.ExpiresIn).Unix(),
    "iat":     time.Now().Unix(),
}
```

### Middleware Flow
1. Extract `Authorization: Bearer <token>` header
2. Validasi format: harus Bearer token
3. Parse & validate JWT dengan secret key
4. Cek signing method (harus HMAC)
5. Set `c.Locals("user_id")` dan `c.Locals("user_email")`
6. Jika invalid, return 401

## 2. Password Security

### Bcrypt Hashing
- **File:** `internal/service/auth_service.go:40`
- Menggunakan `golang.org/x/crypto/bcrypt`
- Default cost (10)
- Password tidak pernah direturn dalam response JSON (`json:"-"` pada field Password)

### Implementation
```go
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// ...
bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
```

## 3. Google reCAPTCHA v2

### Implementation
- **File:** `pkg/captcha/recaptcha.go`
- Memverifikasi token captcha ke Google `siteverify` API
- Dapat dinonaktifkan via `RECAPTCHA_ENABLED=false`
- **Auto-skip** jika `ENV=local` — tidak perlu setup captcha untuk development lokal
- Berlaku untuk endpoint: register, login, forgot-password

### Kapan Captcha di-Skip?

| `ENV` | `RECAPTCHA_ENABLED` | Hasil |
|-------|-------------------|-------|
| `local` | `true` / `false` | **Skip** — cocok untuk local dev |
| `development` | `true` | Validasi captcha |
| `development` | `false` | Skip |
| `production` | `true` | Validasi captcha |
| `production` | `false` | Skip |

### Dua Key yang Dibutuhkan

| Key | Config | Dipakai di | Fungsi |
|-----|--------|-----------|--------|
| **Site Key** | `RECAPTCHA_SITE_KEY` | **Frontend** (HTML/JS) | Untuk render widget "I'm not a robot" |
| **Secret Key** | `RECAPTCHA_SECRET_KEY` | **Backend** (Go) | Untuk verifikasi token ke Google |

> Keduanya didapatkan dari [Google Cloud Console](https://console.cloud.google.com/security/recaptcha) → pilih reCAPTCHA v2 → "I'm not a robot" checkbox.

### Flow
1. **Frontend** pasang widget reCAPTCHA dengan **Site Key**:
   ```html
   <script src="https://www.google.com/recaptcha/api.js" async defer></script>
   <div class="g-recaptcha" data-sitekey="{{ RECAPTCHA_SITE_KEY }}"></div>
   ```
2. User centang "I'm not a robot" → Google generate token
3. **Frontend** kirim token ke backend via field `captcha` di request body
4. **Backend** verifikasi token ke `https://www.google.com/recaptcha/api/v2/siteverify` pakai **Secret Key** (`pkg/captcha/recaptcha.go:37`)
5. Jika invalid, return 422 dengan `CAPTCHA_INVALID`

## 4. Production Security Checklist

### Environment
- [ ] `ENV=production` - Nonaktifkan debug & development mode
- [ ] Ganti semua default secret keys
- [ ] Nonaktifkan database debug logging

### JWT
- [ ] `JWT_SECRET` minimal 32 karakter random
- [ ] Generate menggunakan: `openssl rand -base64 32`
- [ ] Rotasi secret secara berkala

### Database
- [ ] Password database yang kuat
- [ ] Database hanya bisa diakses dari localhost
- [ ] Backup terjadwal

### Network
- [ ] Jangan expose port 3000 ke public
- [ ] Gunakan reverse proxy (Nginx) di port 80/443
- [ ] Aktifkan firewall (UFW)
- [ ] Hanya buka port 22, 80, 443

### Aplikasi
- [ ] Auto-migration hanya dijalankan sekali, lalu di-comment kembali
- [ ] Logging jangan sampai expose data sensitif (password, token)
- [ ] Rate limiting belum diimplementasikan - pertimbangkan untuk menambah
- [ ] Request body size limits via Fiber config

## 5. Potential Improvements

- **Rate Limiting**: Belum ada protection dari brute force login
- **Refresh Token**: Saat ini hanya access token (24 jam), tanpa refresh token
- **CORS**: Belum dikonfigurasi (Fiber default allow all)
- **Input Validation**: Validasi hanya dilakukan di service layer, belum menggunakan validator library
- **SQL Injection**: Aman karena menggunakan GORM parameterized queries
- **XSS/CSRF**: Belum ada protection spesifik
