# Database Schema

## Koneksi
Dikonfigurasi di `config/database/mysql.go` menggunakan GORM dengan driver MySQL.
Database: `domes` (nama database di `.env`: `DB_NAME=domesv2`)

## Tabel: Users (`domes.Users`)

DDL Existing:
```sql
create table domes.Users
(
    id                   int auto_increment primary key,
    username             varchar(255) null unique,
    name                 varchar(255) null,
    first_name           varchar(255) null,
    last_name            varchar(255) null,
    password             varchar(255) null,
    type                 varchar(255) null,
    position             varchar(255) null,
    organization         varchar(255) null,
    phone_number         varchar(255) null,
    email                varchar(255) not null unique,
    registration_id      char(36) null,
    metadata             json null,
    reset_password_token varchar(255) null,
    reset_password_expiry datetime null,
    createdAt            datetime not null,
    updatedAt            datetime not null,
    constraint username unique (username),
    constraint Users_registration_id_foreign_idx
        foreign key (registration_id) references domes.Registrations (id)
) engine = InnoDB;
```

## Go Model (`internal/model/user.go`)

```go
type User struct {
    ID                  uint            `json:"id" gorm:"primaryKey;column:id"`
    Username            *string         `json:"username" gorm:"unique;size:255"`
    Name                *string         `json:"name" gorm:"size:255"`
    FirstName           *string         `json:"first_name" gorm:"column:first_name;size:255"`
    LastName            *string         `json:"last_name" gorm:"column:last_name;size:255"`
    Password            string          `json:"-" gorm:"size:255"`
    Type                *string         `json:"type" gorm:"size:255"`
    Position            *string         `json:"position" gorm:"size:255"`
    Organization        *string         `json:"organization" gorm:"size:255"`
    PhoneNumber         *string         `json:"phone_number" gorm:"column:phone_number;size:255"`
    Email               string          `json:"email" gorm:"unique;size:255;not null"`
    RegistrationID      *string         `json:"registration_id" gorm:"column:registration_id;type:char(36)"`
    Metadata            *string         `json:"metadata" gorm:"type:json"`
    ResetPasswordToken  *string         `json:"-" gorm:"column:reset_password_token;size:255"`
    ResetPasswordExpiry *time.Time      `json:"-" gorm:"column:reset_password_expiry"`
    CreatedAt           time.Time       `json:"created_at" gorm:"column:createdAt"`
    UpdatedAt           time.Time       `json:"updated_at" gorm:"column:updatedAt"`
    DeletedAt           gorm.DeletedAt  `json:"-" gorm:"index"`
}

func (User) TableName() string { return "Users" }
```

## GORM Column Mapping Notes

- Kolom `createdAt` (camelCase) di MySQL -> `gorm:"column:createdAt"` -> tag json `created_at`
- Kolom `updatedAt` (camelCase) di MySQL -> `gorm:"column:updatedAt"` -> tag json `updated_at`
- Kolom `first_name` (snake_case) di MySQL -> `gorm:"column:first_name"`
- `Password` di-exclude dari JSON dengan `json:"-"`
- `DeletedAt` untuk soft delete (gorm.DeletedAt), tidak ada di DDL tapi tidak masalah

## Password Field

- Nilai di database: bcrypt hash dengan prefix `$2a$` (Go) atau `$2b$` (PHP/Node)
- Go's `golang.org/x/crypto/bcrypt` dapat memverifikasi kedua prefix
- Password tidak pernah direturn dalam JSON response

## Request/Response DTOs

### RegisterRequest
```go
type RegisterRequest struct {
    FirstName       string `json:"first_name"`
    LastName        string `json:"last_name"`
    Position        string `json:"position"`
    Organization    string `json:"organization"`
    PhoneNumber     string `json:"phone_number"`
    Email           string `json:"email"`
    Password        string `json:"password"`
    ConfirmPassword string `json:"confirm_password"`
    Captcha         string `json:"captcha,omitempty"`
}
```

### LoginRequest / ForgotPasswordRequest / ResetPasswordRequest
```go
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Captcha  string `json:"captcha,omitempty"`
}

type ForgotPasswordRequest struct {
    Email   string `json:"email"`
    Captcha string `json:"captcha,omitempty"`
}

type ResetPasswordRequest struct {
    Token           string `json:"token"`
    Password        string `json:"password"`
    ConfirmPassword string `json:"confirm_password"`
}
```

## Migration & Seeder

- **Migration**: `database/migrations/001_add_auth_fields.sql`
  - Menambah kolom: first_name, last_name, position, organization, phone_number, reset_password_token, reset_password_expiry
  - Jalan manual: `mysql -u root -p domes < database/migrations/001_add_auth_fields.sql`

- **Seeder**: `database/seeders/001_seed_users.sql`
  - Data admin default + user existing dari CSV
  - Jalan manual: `mysql -u root -p domes < database/seeders/001_seed_users.sql`
