# Development Guide

## Setup

```bash
go mod tidy
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS domes CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p domes < database/migrations/001_add_auth_fields.sql
mysql -u root -p domes < database/seeders/001_seed_users.sql
go run cmd/main.go
```

## Coding Convention

### 1. Package Naming
Lowercase, satu kata: `controller`, `service`, `repository`, `model`

### 2. Interface Pattern
```go
type UserRepository interface {
    Create(user *model.User) error
    FindByEmail(email string) (*model.User, error)
    FindByID(id uint) (*model.User, error)
    FindByResetToken(token string) (*model.User, error)
    Update(user *model.User) error
    UpdatePassword(user *model.User, password string) error
}
```

### 3. Error Handling
- Controller: parse request, delegasikan ke service
- Service: validasi, panggil repository, return AppError
- Repository: wrap database error ke AppError
- Jangan expose error mentah ke client

### 4. Response Format
Gunakan helper dari `pkg/response`:
- `response.Success(c, data, message)` -> 200
- `response.Created(c, data, message)` -> 201
- `response.Error(c, err)` -> dynamic (from AppError)

### 5. Menambah Fitur Baru
1. **Model**: struct di `internal/model/`
2. **Repository**: interface & implementasi di `internal/repository/`
3. **Service**: interface & implementasi di `internal/service/`
4. **Controller**: handler di `internal/controller/`
5. **Routes**: registrasi endpoint + wiring DI di `routes/routes.go`

## Fixes yang sudah dilakukan

### 1. `pkg/response/response.go` — Error field
Sebelumnya `Error` diisi dengan `appErr.Message` (human-readable), bukan error code. Diperbaiki:
```go
// Before (bug)
Error: appErr.Message,
// After (fixed)
Error: appErr.Details,
```

### 2. `internal/model/user.go` — RegisterRequest validation
Typo `validate:"required,min,min=6"` (double min, missing =6). Dihapus karena validasi dilakukan manual di service.

### 3. `internal/controller/auth_controller.go` — Me handler
Sebelumnya hanya return id & email dari JWT claims. Diperbaiki dengan fetch full profile dari database via `AuthService.GetProfile`.

### 4. Model alignment with DDL
- Menambah kolom yang hilang: username, first_name, last_name, type, position, organization, phone_number, metadata, registration_id, reset_password_token, reset_password_expiry
- Menyesuaikan GORM column tags dengan nama kolom di MySQL (createdAt, updatedAt camelCase)
- Menambah `TableName()` return "Users"

### 5. Added missing error code
`INVALID_RESET_TOKEN` di `pkg/errors/errors.go`

## Dependency Injection Changes

`NewAuthService` sekarang menerima `MailService` parameter:
```go
authService := service.NewAuthService(userRepo, mailService)
```
