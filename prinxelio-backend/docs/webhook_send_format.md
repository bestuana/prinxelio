# Format Pengiriman Webhook (API ke n8n)

Data yang dikirim dari API backend ke webhook n8n akan mengikuti format JSON standar ini untuk memastikan pemrosesan yang konsisten. Header `Authorization` menggunakan JWT HS256 ditandatangani dengan `N8N_SECRETKEY_JWT`.

## Struktur Umum

```json
{
  "source": "prinxelio_web",
  "timestamp": 1715000000,
  "action": "NAMA_AKSI",
  "payload": { ... }
}
```

## Detail Kolom

- `source` (string):
  - Sumber pengirim data. Untuk aplikasi ini, nilainya akan selalu `"prinxelio_web"`.

- `timestamp` (integer):
  - Waktu saat data dikirim dalam format Unix timestamp (detik).

- `action` (string):
  - Menjelaskan jenis aksi yang harus dilakukan oleh n8n. Nilai ini akan bervariasi tergantung pada tujuan webhook.
  - Contoh Aksi:
    - `SEND_OTP`: Untuk mengirimkan kode OTP ke nomor WhatsApp pengguna.
    - `CREATE_TRANSACTION`: Untuk membuat transaksi baru di payment gateway.

- `payload` (object):
  - Berisi data utama yang dibutuhkan oleh n8n untuk melakukan aksi tersebut.

## Contoh Pengiriman

### 1. Mengirim OTP

- **Action**: `SEND_OTP`
- **Payload**: Memuat nomor telepon target dan kode OTP yang akan dikirim.

```json
{
  "source": "prinxelio_web",
  "timestamp": 1678886400,
  "action": "SEND_OTP",
  "payload": {
    "phone": "6281234567890",
    "otp": "123456"
  }
}
```

Header:

```
Authorization: Bearer <jwt-signed-with-N8N_SECRETKEY_JWT>
Content-Type: application/json
User-Agent: Prinxelio-Backend/1.0
```

### 2. Membuat Transaksi Pembayaran

- **Action**: `CREATE_TRANSACTION`
- **Payload**: Memuat detail transaksi seperti ID referensi, jumlah pembayaran, dan detail produk.

```json
{
  "source": "prinxelio_web",
  "timestamp": 1678886500,
  "action": "CREATE_TRANSACTION",
  "payload": {
    "reference_id": "TRX-2023-ABCDE",
    "amount": 125000,
    "product_name": "Ebook Panduan Go Lengkap",
    "customer_phone": "6281234567890"
  }
}
```

Header sama seperti pengiriman OTP.
