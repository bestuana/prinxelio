# Dokumentasi Kategori

## Perubahan Skema Database
- Tabel `category` ditambah kolom:
  - `category_images` (TEXT) untuk URL gambar kategori
  - `category_color` (VARCHAR(7)) untuk kode warna hex

## Migrasi
- Migrasi diterapkan otomatis saat server start melalui `database.InitDB()`.
- Jika kolom belum ada, server menambahkannya via `ALTER TABLE`.

## Seed Data
Contoh minimal 5 record dimasukkan otomatis saat start:
```json
{
  "id": 1,
  "category_name": "Website",
  "category_create_at": "2023-01-01",
  "category_images": "https://example.com/cat1.jpg",
  "category_color": "#FF5733"
}
```
Kategori lain: `Ebook`, `Desain`, `AI Tools`, `Audio`.

## API Kategori
- `GET /api/categories?q=<optional>`
  - Response:
  ```json
  {
    "status": true,
    "message": "Berhasil",
    "data": [
      {"id":1,"category_name":"Website","category_images":"...","category_color":"#FF5733"}
    ]
  }
  ```

## UI Slider Kategori
- Section: `Daftar Kategori` di atas `Daftar Produk`.
- Slider horizontal dengan item persegi (1:1) berisi gambar, overlay warna, dan nama.
- Posisi default slider di tengah.
- Klik kategori akan memberi outline aktif dan memfilter produk.

## Alur Interaksi Kategori-Produk
- Pengguna klik kategori → state `selectedCategory` diset → grid produk dirender ulang memfilter `category_name` yang sama.
- Klik kategori yang sama: filter tetap.
- Klik kategori berbeda: filter diperbarui ke kategori baru.

## Contoh Payload
- Request produk: `GET /api/products`
- Response produk (contoh item):
```json
{"id":2,"product_name":"Template Website Portofolio","product_image":"...","product_price":75000,"product_discount":5000,"product_discount_amount":20,"product_desc":"...","product_viewed":0,"product_downloaded":0,"category_name":"Website"}
```

