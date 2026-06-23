# Overview Project

## Nama Project
**DOMESv2 Backend API** - Sebuah REST API backend modern berbasis GoFiber.

## Tujuan
Project ini adalah backend service untuk aplikasi "DOMES" yang menyediakan:
- Autentikasi pengguna (register/login) dengan JWT
- Manajemen profil pengguna (first_name, last_name, position, organization, phone_number)
- Forgot & reset password
- Health checking untuk monitoring sistem
- Perlindungan endpoint dengan Google reCAPTCHA v2

## Tech Stack

| Komponen | Teknologi | Versi |
|----------|-----------|-------|
| Bahasa | Go | 1.21 |
| Web Framework | GoFiber v2 | 2.52.0 |
| ORM | GORM | 1.25.12 |
| Database | MySQL | 8.0+ (tabel: `domes.Users`) |
| Cache | Redis (opsional) | 6.0+ |
| JWT | golang-jwt | v5 |
| Logging | Uber Zap | 1.28.0 |
| Captcha | Google reCAPTCHA v2 | - |
| Password | bcrypt (`$2a$` / `$2b$`) | - |

## Status Project
- **Initial commit** - Project masih tahap awal
- Database existing dengan tabel `Users` yang memiliki data riil
- Auto-migration tidak digunakan (manual SQL migration)
- Seeder tersedia untuk data awal

## Struktur Direktori
```
DOMESV2-GOFIBER/
  .env                  # Environment variables
  go.mod / go.sum       # Go modules
  cmd/
    main.go             # Entrypoint aplikasi
  config/
    env.go              # Config loader
    database/mysql.go   # Koneksi MySQL/GORM
    redis/redis.go      # Koneksi Redis
    logger/logger.go    # Inisialisasi Zap logger
  database/
    migrations/         # SQL migration (manual)
    seeders/            # SQL seeder (manual)
  routes/
    routes.go           # Definisi route dan DI wiring
  internal/
    model/              # Struct models & DTOs
    repository/         # Layer data (GORM queries)
    service/            # Layer business logic
    controller/         # Layer HTTP handlers
    middleware/          # JWT & logging middleware
  pkg/
    captcha/            # Google reCAPTCHA client
    errors/             # Custom error types
    response/           # Standardized API response
  docsagent/            # Dokumentasi internal AI agent
  API-REFERENCE.md      # Dokumentasi API untuk user
```
