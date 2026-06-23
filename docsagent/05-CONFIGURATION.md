# Configuration Reference

Semua konfigurasi di-load dari file `.env` atau system environment variables.

## File `.env`

```
# =====================
# Server Configuration
# =====================
PORT=3000
ENV=development

# =====================
# MySQL Database
# =====================
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=domesv2
DB_CHARSET=utf8mb4
DB_PARSE_TIME=True
DB_LOC=Local

# =====================
# Redis (Optional)
# =====================
REDIS_ENABLED=false
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# =====================
# JWT
# =====================
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=24h

# =====================
# Google reCAPTCHA v2
# =====================
RECAPTCHA_SECRET_KEY=your-recaptcha-secret-key
RECAPTCHA_SITE_KEY=your-recaptcha-site-key
RECAPTCHA_ENABLED=true
```

## Tabel Konfigurasi

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| **Server** | | |
| `PORT` | `3000` | Port HTTP server |
| `ENV` | `development` | Environment mode (`local`/`development`/`production`)<br>`local` = skip captcha, debug mode |
| **Database** | | |
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASSWORD` | `""` | MySQL password |
| `DB_NAME` | `domesv2` | MySQL database name |
| `DB_CHARSET` | `utf8mb4` | Character set |
| `DB_PARSE_TIME` | `true` | Parse time.Time from MySQL |
| `DB_LOC` | `Local` | Timezone location |
| **Redis** | | |
| `REDIS_ENABLED` | `false` | Aktifkan/nonaktifkan Redis |
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_PASSWORD` | `""` | Redis password |
| `REDIS_DB` | `0` | Redis database number |
| **JWT** | | |
| `JWT_SECRET` | `your-super-secret-jwt-key` | Secret key untuk signing JWT |
| `JWT_EXPIRES_IN` | `24h` | Masa berlaku token (format Go duration) |
| **Captcha** | | |
| `RECAPTCHA_SECRET_KEY` | `""` | Secret key dari Google Cloud Console |
| `RECAPTCHA_SITE_KEY` | `""` | Site key dari Google Cloud Console |
| `RECAPTCHA_ENABLED` | `true` | Aktifkan/nonaktifkan captcha |

## Environment Mode Behavior

| `ENV` | Captcha | Logging | Notes |
|-------|---------|---------|-------|
| `local` | **Skipped** | Development (console) | Untuk development lokal tanpa perlu setup captcha |
| `development` | Depend on `RECAPTCHA_ENABLED` | Development (console) | Development dengan captcha aktif |
| `production` | Depend on `RECAPTCHA_ENABLED` | Production (file JSON) | Production |

## Production Checklist

- [ ] `ENV=production` - Nonaktifkan debug mode
- [ ] `JWT_SECRET` - Ganti dengan random string min 32 karakter
- [ ] `DB_PASSWORD` - Password database yang kuat
- [ ] `RECAPTCHA_SECRET_KEY` - Secret key dari Google Cloud Console
- [ ] `RECAPTCHA_ENABLED=true` - Pastikan captcha aktif
- [ ] `PORT` - Jangan expose ke public (gunakan reverse proxy)
