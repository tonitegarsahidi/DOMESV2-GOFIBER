# Error Handling

## AppError Struct

**File:** `pkg/errors/errors.go`

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

## Error Constructors

| Constructor | HTTP Status | Contoh Penggunaan |
|-------------|-------------|-------------------|
| `NewBadRequestError` | 400 | Invalid request body, invalid reset token |
| `NewUnauthorizedError` | 401 | Invalid credentials |
| `NewForbiddenError` | 403 | Access denied |
| `NewNotFoundError` | 404 | User not found |
| `NewConflictError` | 409 | Duplicate email |
| `NewInternalServerError` | 500 | Database error, hash error |
| `NewValidationError` | 422 | Invalid captcha, validasi input gagal |

## Error Codes

```go
const (
    ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
    ErrCodeTokenExpired       = "TOKEN_EXPIRED"
    ErrCodeInvalidToken       = "INVALID_TOKEN"
    ErrCodeUserExists         = "USER_ALREADY_EXISTS"
    ErrCodeUserNotFound       = "USER_NOT_FOUND"
    ErrCodeCaptchaInvalid     = "CAPTCHA_INVALID"
    ErrCodeCaptchaMissing     = "CAPTCHA_MISSING"
    ErrCodeDatabaseError      = "DATABASE_ERROR"
    ErrCodeRedisError         = "REDIS_ERROR"
    ErrCodeInternalServer     = "INTERNAL_ERROR"
    ErrCodeValidationFailed   = "VALIDATION_FAILED"
    ErrCodeInvalidResetToken  = "INVALID_RESET_TOKEN"
)
```

## Response Format

### Error Response (fixed)
```json
{
  "success": false,
  "message": "Human-readable error message",
  "error": "ERROR_CODE",           // Machine-readable code (FIXED)
  "details": "message: ERROR_CODE"   // Full details
}
```

## Error Flow

1. **Repository**: mapping error database ke AppError
   - Duplicate entry -> `NewConflictError("...", "USER_ALREADY_EXISTS")`
   - Record not found -> `NewNotFoundError("...", "USER_NOT_FOUND")`
   - Invalid/expired reset token -> `NewBadRequestError("...", "INVALID_RESET_TOKEN")`

2. **Service**: validasi business logic
   - Empty fields -> `NewValidationError("...", "VALIDATION_FAILED")`
   - Password mismatch -> `NewValidationError("...", "VALIDATION_FAILED")`
   - Wrong password -> `NewUnauthorizedError("...", "INVALID_CREDENTIALS")`

3. **Response.Error()**: 
   - Jika `*AppError`: map `Code` -> HTTP status, `Details` -> `error` field, `Error()` -> `details` field
   - Jika bukan AppError: return 500 Internal Server Error
