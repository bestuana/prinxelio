# Format Respons API (API ke Frontend)

Semua output API dari backend ke frontend akan mengikuti format JSON standar ini untuk memastikan konsistensi.

## Struktur Umum

```json
{
  "status": true,
  "message": "Pesan yang mendeskripsikan hasil operasi",
  "data": { ... } 
}
```

## Detail Kolom

- `status` (boolean):
  - `true`: Menandakan bahwa permintaan berhasil diproses.
  - `false`: Menandakan bahwa terjadi kesalahan atau kegagalan dalam memproses permintaan.

- `message` (string):
  - Memberikan pesan yang jelas dan ringkas mengenai hasil dari permintaan. Contoh: "Berhasil memuat data produk", "Nomor WhatsApp tidak valid", atau "Terjadi kesalahan internal server".

- `data` (object | array | null):
  - Berisi data utama yang diminta oleh klien. Strukturnya bisa berupa objek tunggal, sebuah array dari beberapa objek, atau `null` jika tidak ada data yang dikembalikan (misalnya pada operasi penghapusan).

## Contoh Sukses

Permintaan untuk mendapatkan detail produk:

```json
{
  "status": true,
  "message": "Berhasil memuat data produk",
  "data": {
    "id": 1,
    "product_name": "Ebook Panduan Go Lengkap",
    "product_image": "https://example.com/image.png",
    "product_price": 150000,
    "product_discount": 120000,
    "product_discount_amount": 20,
    "product_desc": "Panduan lengkap belajar Go."
  }
}
```

## Contoh Gagal

Permintaan dengan parameter yang tidak valid:

```json
{
  "status": false,
  "message": "ID produk tidak ditemukan.",
  "data": null
}
```
