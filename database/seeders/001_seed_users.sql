-- Seeder: 001_seed_users.sql
-- Description: Data awal untuk tabel Users
--
-- Cara menjalankan:
--   mysql -u root -p domes < database/seeders/001_seed_users.sql
--
-- Note: Password di bawah ini adalah bcrypt hash dari "admin123"
-- Gunakan $2a$ (format Go) atau $2b$ (format PHP/Node) - keduanya kompatibel

USE domes;

INSERT INTO Users (username, name, first_name, last_name, password, type, position, organization, phone_number, email, createdAt, updatedAt)
VALUES
    ('admin', 'Admin Domes', 'Admin', 'Domes', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', 'System Administrator', 'UNITED NATIONS', '+62123456789', 'admin@domes.un.org', NOW(), NOW()),
    ('erlangga_admin', 'Erlangga Agustino Landiyanto', 'Erlangga', 'Agustino Landiyanto', '$2b$10$JUibY3K9bG10MdzcAh6zXOr.5SdO8knCT0KqDicdRRBqMejyqSOQG', 'admin', 'Administrator', 'UNITED NATIONS', '+628123456789', 'erlangga.landiyanto@un.org', '2023-06-08 01:03:49', '2024-05-13 04:01:50')
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    first_name = VALUES(first_name),
    last_name = VALUES(last_name),
    type = VALUES(type);
