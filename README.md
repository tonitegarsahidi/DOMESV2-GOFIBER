# DOMESv2 Backend - GoFiber

Backend API modern berbasis Go-Fiber dengan MySQL dan Redis (opsional). Dibangun dengan prinsip separation of concerns dan best practices development.

## 🛠️ Tech Stack
- **Framework**: GoFiber v2 (fasthttp-based web framework)
- **Database**: MySQL dengan GORM ORM
- **Cache**: Redis (opsional)
- **Auth**: JWT
- **Captcha**: Google reCAPTCHA v2
- **Logging**: Zap (high-performance logging)
- **Env Management**: godotenv

## 📋 Prasyarat
- Go 1.21 atau lebih baru
- MySQL 8.0+
- Redis 6.0+ (opsional)
- Git

## 🚀 Instalasi & Menjalankan di Local Development

### 1. Clone Repository
```bash
git clone <repository-url>
cd DOMESV2-GOFIBER
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Setup Environment Variables
Salin file `.env` dan sesuaikan dengan konfigurasi lokal Anda:
```bash
# File .env sudah tersedia, tinggal edit nilai-nilai berikut:
nano .env
```

Konfigurasi penting yang perlu diubah:
```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password_database_anda
DB_NAME=domesv2

# Redis (opsional, set true jika ingin pakai Redis)
REDIS_ENABLED=false

# JWT (ganti dengan secret key yang aman di production!)
JWT_SECRET=your-super-secret-jwt-key-change-in-production-32charsmin

# reCAPTCHA (dapatkan dari Google Cloud Console)
RECAPTCHA_SECRET_KEY=your-recaptcha-secret-key
RECAPTCHA_SITE_KEY=your-recaptcha-site-key
```

### 4. Buat Database MySQL
Login ke MySQL dan buat database:
```sql
CREATE DATABASE domesv2 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. Aktifkan Auto Migration
Buka file `cmd/main.go`, uncomment baris berikut untuk auto-create tabel:
```go
// Uncomment ini untuk pertama kali running
db := database.GetDB()
db.AutoMigrate(&model.User{})
```

### 6. Jalankan Aplikasi
```bash
# Development mode
go run cmd/main.go
```

Server akan berjalan di `http://localhost:3000`

### 7. Test API
Cek health check untuk memastikan semua berjalan:
```bash
curl http://localhost:3000/api/health-check
```

## 📡 API Endpoints

Backend ini menyediakan endpoint lengkap untuk sistem manajemen dokumen PBB:
* **Authentication & Profiles:** Registrasi, login, edit profil, ganti password, preferensi notifikasi.
* **Admin Whitelist Settings:** Whitelist email admin.
* **Reference Data:** Data SDGs, PBB Agencies, Sectors, Languages, Joint Programmes, dll.
* **Public Documents & Search:** Pencarian dokumen, list dokumen, detail, related docs, tracking download.
* **Broken Link Reports:** Pelaporan tautan rusak oleh publik dan manajemen status laporan.
* **CMS Dashboard & Submissions:** Draft submissions wizard (Step 1-4), publishing/unpublishing dokumen.
* **CMS User Management:** CRUD akun pengelola (admin/editor) oleh administrator.
* **File Upload:** Upload file PDF/Word, cover dokumen, avatar user, dan validasi tautan eksternal.

Dokumentasi lengkap kontrak API: [API-REFERENCE.md](API-REFERENCE.md)

## 🗄️ Database Migration & Seeder

Database migration dan seeding data referensi (seperti SDGs, Agencies, Sectors, LNOBs, dll.) kini berjalan **secara otomatis (Auto-run)** pada saat aplikasi pertama kali dijalankan. Aplikasi mendeteksi perubahan skema model dan melakukan pengisian data awal secara dinamis.

Jika Anda ingin melakukan manipulasi database manual atau melihat struktur awal:
* File seeder dan migrasi opsional berada di folder `database/`.
* Pastikan kredensial database di file `.env` sudah sesuai sebelum menjalankan aplikasi.

## 🌐 Deployment ke Production Server

### 1. Build Binary untuk Production
Compile aplikasi menjadi binary file:
```bash
go build -o domesv2 ./cmd/main.go
```

### 2. Upload ke Server
Upload binary, .env, dan pastikan folder `logs` bisa dibuat oleh aplikasi:
```bash
# Buat folder logs di server
mkdir -p /opt/domesv2/logs
chmod 755 /opt/domesv2/logs
```

### 3. Setup Systemd Service (untuk auto-restart)
Buat file service systemd:
```bash
sudo nano /etc/systemd/system/domesv2.service
```

Isi dengan konfigurasi berikut:
```ini
[Unit]
Description=DOMESv2 Backend API
After=network.target mysql.service redis-server.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/domesv2
Environment=ENV=production
ExecStart=/opt/domesv2/domesv2
Restart=always
RestartSec=5
StandardOutput=journal+console
StandardError=journal+console

[Install]
WantedBy=multi-user.target
```

### 4. Jalankan Service
```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service agar jalan saat boot
sudo systemctl enable domesv2

# Start service
sudo systemctl start domesv2

# Cek status
sudo systemctl status domesv2
```

### 5. Setup Reverse Proxy dengan Nginx
Install Nginx dan buat konfigurasi:
```bash
sudo nano /etc/nginx/sites-available/domesv2.yourdomain.com
```

Konfigurasi Nginx:
```nginx
server {
    listen 80;
    server_name domesv2.yourdomain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Aktifkan konfigurasi:
```bash
sudo ln -s /etc/nginx/sites-available/domesv2.yourdomain.com /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 6. SSL dengan Let's Encrypt (Certbot)
```bash
sudo certbot --nginx -d domesv2.yourdomain.com
```

### 7. Monitoring Logs
Cek logs aplikasi:
```bash
# Log dari systemd
journalctl -u domesv2 -f

# Log file aplikasi
tail -f /opt/domesv2/logs/app.log
```

## 🔧 Maintenance & Troubleshooting

### Restart Service
```bash
sudo systemctl restart domesv2
```

### Cek Logs Error
```bash
journalctl -u domesv2 --since "1 hour ago"
```

### Update Aplikasi
```bash
# Pull kode baru
git pull

# Rebuild
go build -o domesv2 ./cmd/main.go

# Restart service
sudo systemctl restart domesv2
```

## 🔒 Keamanan Production Checklist
- [ ] Ubah `JWT_SECRET` dengan nilai yang sangat kuat (min 32 karakter random)
- [ ] Set `ENV=production` di .env
- [ ] Nonaktifkan debug mode
- [ ] Setup firewall (ufw) hanya buka port 80, 443, 22
- [ ] Database password yang kuat
- [ ] Jangan expose port 3000 ke public
- [ ] Aktifkan UFW firewall
- [ ] Backup database secara teratur

## 📝 Lisensi
MIT License