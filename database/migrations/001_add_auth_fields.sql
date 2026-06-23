-- Migration: 001_add_auth_fields.sql
-- Description: Menambahkan kolom untuk fitur register, forgot-password, dan profile
-- 
-- Cara menjalankan:
--   mysql -u root -p domes < database/migrations/001_add_auth_fields.sql

USE domes;

ALTER TABLE Users
    ADD COLUMN first_name VARCHAR(255) NULL AFTER name,
    ADD COLUMN last_name VARCHAR(255) NULL AFTER first_name,
    ADD COLUMN position VARCHAR(255) NULL AFTER type,
    ADD COLUMN organization VARCHAR(255) NULL AFTER position,
    ADD COLUMN phone_number VARCHAR(255) NULL AFTER organization,
    ADD COLUMN reset_password_token VARCHAR(255) NULL,
    ADD COLUMN reset_password_expiry DATETIME NULL;
