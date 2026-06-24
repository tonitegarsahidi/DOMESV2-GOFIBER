# Arsitektur Project

## Layered Architecture

Project menggunakan arsitektur berlapis (Controller -> Service -> Repository):

```
Request HTTP
     |
     v
  [Middleware] (Logging, JWT, Recover)
     |
     v
  [Controller]  -> HTTP handler, parse request, kirim response
     |
     v
  [Service]     -> Business logic, validasi, orchestrasi
     |
     v
  [Repository]  -> Akses data (GORM ke MySQL), Redis
     |
     v
  [Database/Redis]
```

## Dependency Injection (Manual)

Semua dependensi di-wiring secara manual di `routes/routes.go:12-23`:

```go
userRepo := repository.NewUserRepository()
healthRepo := repository.NewHealthRepository()
mailService := service.NewMailService()

authService := service.NewAuthService(userRepo, mailService)
healthService := service.NewHealthService(healthRepo)

authController := controller.NewAuthController(authService)
healthController := controller.NewHealthController(healthService)
```

## Interface-Based Abstraksi

Setiap layer menggunakan **Go interfaces** untuk testability:

- `repository.UserRepository` - `Create`, `FindByEmail`, `FindByID`, `FindByResetToken`, `Update`, `UpdatePassword`
- `repository.HealthRepository` - `CheckDatabase`, `CheckRedis`
- `service.AuthService` - `Register`, `Login`, `ForgotPassword`, `ResetPassword`, `GetProfile`
- `service.HealthService` - `CheckHealth`
- `service.MailService` - `SendResetPassword`

## Alur Autentikasi

### Register
1. Client POST `/api/v2/auth/register` dengan `{first_name, last_name, position, organization, phone_number, email, password, confirm_password, captcha}`
2. `AuthController.Register` -> parse body
3. `AuthService.Register` -> verify captcha -> validasi input -> hash password (bcrypt) -> auto-generate username dari email -> `UserRepository.Create` -> generate JWT
4. Response 201: `{token, user}`

### Login
1. Client POST `/api/v2/auth/login` dengan `{email, password, captcha}`
2. `AuthController.Login` -> parse body
3. `AuthService.Login` -> verify captcha -> `UserRepository.FindByEmail` -> compare bcrypt -> generate JWT
4. Response 200: `{token, user}`

### Forgot Password
1. Client POST `/api/v2/auth/forgot-password` dengan `{email, captcha}`
2. `AuthService.ForgotPassword` -> verify captcha -> cari user -> generate random token (32 byte hex) -> simpan token + expiry (1 jam) -> kirim email
3. Response 200: pesan sukses (tetap sama meski email tidak ditemukan)

### Reset Password
1. Client POST `/api/v2/auth/reset-password` dengan `{token, password, confirm_password}`
2. `AuthService.ResetPassword` -> validasi input -> cari user by token (cek expiry) -> hash password baru -> update password, hapus token

### JWT Validation
1. Client GET `/api/v2/user/me` dengan header `Authorization: Bearer <token>`
2. `JWTMiddleware` -> extract token -> parse & validate -> set `c.Locals("user_id")` dan `c.Locals("user_email")`
3. `AuthController.Me` -> panggil `AuthService.GetProfile` yang fetch dari DB -> return full profile

## Database

Tabel `domes.Users` sudah ada dengan struktur lengkap:

| Kolom | Tipe | Notes |
|-------|------|-------|
| id | int (PK, auto_increment) | |
| username | varchar(255) UNIQUE | Auto-generate dari email |
| name | varchar(255) | Auto-fill: "First Last" |
| first_name | varchar(255) | Dari register |
| last_name | varchar(255) | Dari register |
| password | varchar(255) | bcrypt hash |
| type | varchar(255) | admin/user |
| position | varchar(255) | Jabatan |
| organization | varchar(255) | Organisasi |
| phone_number | varchar(255) | No telepon |
| email | varchar(255) UNIQUE NOT NULL | |
| registration_id | char(36) FK -> Registrations | |
| metadata | json | |
| reset_password_token | varchar(255) | Token reset |
| reset_password_expiry | datetime | Expiry token |
| createdAt | datetime | GORM column:createdAt |
| updatedAt | datetime | GORM column:updatedAt |

## Migration & Seeder

- Migration: `database/migrations/001_add_auth_fields.sql` - menambah kolom baru ke tabel Users
- Seeder: `database/seeders/001_seed_users.sql` - data awal admin
- Keduanya dijalankan manual via `mysql` CLI, tidak auto-run oleh Go
