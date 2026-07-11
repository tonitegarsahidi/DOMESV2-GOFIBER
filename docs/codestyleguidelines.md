# Backend Code Style Guidelines — DOMES V2

Dokumen ini mendefinisikan standar penulisan kode, manajemen error, database, serta pedoman pengembangan di bagian Backend **DOMES V2** (Go Fiber).

---

## 🐹 Go Coding Style & Formatting

1. **Format Kode Standard**:
   * Selalu jalankan `gofmt` dan `goimports` sebelum melakukan commit.
   * Panjang baris kode sebaiknya tidak melebihi 120 karakter untuk keterbacaan yang optimal.
2. **Penamaan Pengenal (Naming Conventions)**:
   * **CamelCase**: Gunakan PascalCase untuk identifier yang diekspor (fungsi, struct, interface yang berawalan huruf besar) dan camelCase untuk yang internal/lokal.
   * Singkatan umum harus ditulis secara konsisten dalam huruf besar seluruhnya (contoh: `APIURL` atau `DocumentUUID`, bukan `ApiUrl` atau `DocumentUuid`).
   * Nama file ditulis dengan format `snake_case.go` (contoh: `document_repository.go`, `jwt_middleware.go`).

---

## 🏗️ Layered Architecture Rules

Setiap layer memiliki tanggung jawab khusus yang tidak boleh dilompati:
1. **Controller Layer**:
   * Hanya bertanggung jawab memproses input HTTP request, memanggil service yang sesuai, dan mengembalikan format standard response.
   * **Dilarang keras** melakukan query database langsung atau mengolah logika bisnis kompleks di Controller.
2. **Service Layer**:
   * Tempat logika bisnis utama ditulis. Menerima data ter-parse dari controller, memproses transaksi, dan memanggil repository.
   * Tidak boleh tahu-menahu tentang konteks HTTP (`*fiber.Ctx` tidak boleh dilewatkan ke Service).
3. **Repository Layer**:
   * Berisi query SQL/GORM. Hanya berinteraksi dengan database dan mengembalikan model data ke Service.

---

## 📝 Logging & Observability

1. **Zap Logger**:
   * Selalu gunakan Zap Logger (`pkg/logger` atau `zap.L()`) untuk mencatat log aktivitas sistem dan error.
   * **Dilarang** menggunakan `fmt.Println()`, `print()`, atau library log bawaan Go (`log.Println`) di dalam controller, service, maupun repository.
2. **Log Levels**:
   * **Info**: Untuk mencatat alur aplikasi normal (contoh: startup, cron job run, database connected).
   * **Warn**: Untuk kondisi yang tidak diinginkan tapi tidak menghentikan aplikasi (contoh: request bad input, captcha verification failed).
   * **Error**: Untuk kegagalan sistem internal, error query database, kegagalan pihak ketiga (email service down).

---

## 🚨 Error Handling & Standard Responses

1. **Centralized Error Handler**:
   * Semua error harus dikembalikan ke atas (*bubble up*) hingga Controller, lalu ditangani oleh middleware terpusat di `internal/middleware/error_handler.go`.
   * Gunakan helper di `pkg/errors` untuk mendefinisikan tipe error (seperti *NotFound*, *Unauthorized*, *BadRequest*, *InternalServerError*).
2. **Consistent JSON Response**:
   * Semua response API wajib dibungkus menggunakan struct standard dari `pkg/response` dengan format:
     ```json
     {
       "success": true,
       "message": "Pesan sukses",
       "data": { ... }
     }
     ```
     Atau jika terjadi error:
     ```json
     {
       "success": false,
       "message": "Pesan kesalahan ramah pengguna",
       "error": "Detail error teknis (hanya tampil jika mode development)"
     }
     ```

---

## 🗄️ Database & GORM Best Practices

1. **Gunakan Parameterized Queries**:
   * Selalu gunakan fitur binding parameter dari GORM untuk menghindari SQL Injection. Jangan menggabungkan string mentah ke dalam query SQL.
2. **Auto-Migrations**:
   * Skema tabel baru wajib didaftarkan di fungsi AutoMigrate `cmd/main.go` / `cmd/migrate/main.go`.
   * Pastikan nama tabel mengikuti penamaan versi 2 (awalan `V2` dan bentuk jamak, contoh: `V2Documents`).
3. **Isolation in Tests**:
   * Pengujian unit (`go test`) tidak boleh merusak atau menghapus data operational database utama. Pastikan koneksi DB test diarahkan ke database test terpisah atau menggunakan transaksi yang di-rollback.
