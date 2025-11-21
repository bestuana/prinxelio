# Panduan Setup dan Deployment

## Prasyarat
- Go 1.21+
- MySQL database dengan kredensial sesuai `.env`
- Akses ke n8n dan aktivasi workflow untuk OTP dan transaksi

## Konfigurasi
1. Salin file `configs/.env` dan sesuaikan nilai sesuai environment produksi.
2. Pastikan `N8N_WEBHOOK_URL_OTP` dan `N8N_WEBHOOK_URL_TRANSACTION_CREATE` aktif. Gunakan metode `POST`.

## Menjalankan secara lokal
```bash
cd prinxelio-backend
go run cmd/server/main.go
```
Server berjalan di `http://localhost:8080`.

## Build dan Deploy
```bash
go build -o bin/prinxelio cmd/server/main.go
```
Jalankan binary:
```bash
./bin/prinxelio
```
Gunakan process manager (systemd, pm2, atau supervisor) untuk menjaga proses tetap berjalan.

## Keamanan
- Jangan pernah mengekspos isi `.env` di repositori publik.
- Header `Authorization: Bearer` untuk webhook n8n harus menggunakan `N8N_SECRETKEY_JWT`.
- JWT untuk klien ditandatangani dengan `PUBLIC_API_SECRET_KEY` dengan durasi `TIME_JWT`.

## Migrasi Database
Inisialisasi tabel dilakukan otomatis saat server start melalui `database.InitDB()`.

## Postman dan OpenAPI
- Import `docs/postman_collection.json` ke Postman.
- Gunakan `docs/openapi.yaml` untuk integrasi dan dokumentasi.

