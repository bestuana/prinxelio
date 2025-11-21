# Format Respons Webhook (n8n ke API)

Respons yang diharapkan dari n8n untuk endpoint terkait transaksi dan OTP.

## 1. Respons Pengiriman OTP

Endpoint backend: `/api/otp/send` mengirim envelope `SEND_OTP` ke `N8N_WEBHOOK_URL_OTP`.

Respons yang diharapkan:

```json
{
  "status": true,
  "data": {
    "delivered": true,
    "message": "Berhasil dikirimkan!"
  }
}
```

atau

```json
{
  "status": false,
  "data": {
    "delivered": false,
    "message": "Gagal dikirimakn!"
  }
}
```

Backend akan meneruskan nilai `status` dan `message` dari respons webhook ke frontend. Jika webhook mengembalikan `status=false` atau `delivered=false`, API akan merespons dengan `status=false` dan `data.delivered` sesuai.

## 2. Respons Pembuatan Transaksi

Endpoint backend: `/api/transactions/create` mengirim envelope `CREATE_TRANSACTION` tanpa `reference` ke `N8N_WEBHOOK_URL_TRANSACTION_CREATE`.

Payload terkirim:

```json
{
  "amount": 120000,
  "customer_phone": "6281234567890",
  "product_name": "Ebook Panduan Go Lengkap"
}
```

Respons yang diharapkan:

```json
{
  "success": true,
  "data": {
    "reference": "DEV-T39871310233VMC74",
    "qr_url": "https://tripay.co.id/qr/DEV-T39871310233VMC74",
    "merchant_ref": "INV-1763612693105",
    "customer_name": "USER-...",
    "customer_email": "USER-...@mail.com",
    "customer_phone": "6281234567890",
    "amount": 121590,
    "total_fee": 1590,
    "status": "UNPAID",
    "expired_time": 1763613473
  }
}
```

Jika `success=false`, backend mengembalikan pesan gagal tanpa mencatat transaksi.
Backend akan meneruskan `success` dan `data.message` dari webhook ke frontend pada respons API.

## 3. Webhook Status Pembayaran

n8n memanggil backend pada `/api/webhook/payment-status` dengan header `Authorization: Bearer N8N_SECRETKEY_JWT`.

Payload:

```json
{
  "reference": "TRX-1699999999-123456",
  "status": "PAID"
}
```

Nilai `status` valid: `PENDING`, `PAID`, `FAILED`, `EXPIRED`.

Backend akan mengupdate tabel `transactions.status` dan melakukan broadcast ke klien WebSocket/SSE yang terhubung pada reference tersebut.
