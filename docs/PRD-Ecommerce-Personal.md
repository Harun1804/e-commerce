# Product Requirements Document (PRD)
# Platform E-Commerce Pribadi (Produk & Jasa)

| | |
|---|---|
| **Versi** | 1.0 |
| **Tanggal** | 7 Juli 2026 |
| **Status** | Draft |
| **Pemilik Produk** | Personal Project |
| **Status Tech Stack** | **Final** (lihat Bab 6) |

---

## 1. Latar Belakang

Dibutuhkan sebuah platform e-commerce pribadi yang dapat menjual **dua jenis komoditas sekaligus**:

1. **Produk fisik/digital** — dengan stok, varian, dan pengiriman.
2. **Jasa/layanan** — berbasis booking/jadwal, tanpa stok fisik, dengan durasi dan slot waktu.

Platform dibangun dengan arsitektur **terpisah antara Backend (BE) dan Frontend (FE)** (headless/decoupled), agar lebih fleksibel untuk pengembangan, scaling, dan potensi ekspansi ke mobile app di masa depan. Penyimpanan file (gambar produk, dokumen, bukti transfer, dsb.) menggunakan **MinIO (S3-compatible object storage)**.

Proyek ini direncanakan untuk **jangka panjang (long-term)**. Saat ini dijalankan sebagai **single-tenant** (satu toko pribadi), namun seluruh arsitektur (database, auth, storage, struktur kode) **dirancang agar siap direfactor menjadi SaaS multi-tenant** di masa depan tanpa perlu menulis ulang dari nol (lihat Bab 13).

---

## 2. Tujuan (Goals)

- Membangun platform jualan online yang mendukung model **produk** dan **jasa** dalam satu sistem.
- Memisahkan BE dan FE agar arsitektur scalable dan maintainable.
- Menyediakan sistem pembayaran (sampai status **paid**), manajemen pesanan, dan manajemen booking jasa.
- Menyediakan dashboard admin untuk mengelola produk, jasa, pesanan, dan pengguna.
- Menggunakan object storage (MinIO) yang self-hosted untuk efisiensi biaya dibanding cloud storage berbayar.
- Menyiapkan pondasi arsitektur agar dapat berkembang menjadi **produk SaaS multi-tenant** tanpa migrasi besar di kemudian hari.

### Cakupan v1 (Saat Ini)
- Alur order **berhenti setelah pembayaran berhasil** (status `paid`). Proses pengiriman fisik **belum** ditangani sistem di v1 — masih manual/di luar sistem.

### Non-Goals (di luar cakupan v1)
- Modul **pengiriman/delivery** dan **manajemen armada** — direncanakan sebagai **Fase Lanjutan** (lihat Bab 11, Fase 5).
- Multi-tenant aktif (pendaftaran seller lain) — arsitektur disiapkan, tapi **tidak diaktifkan** di v1.
- Aplikasi mobile native (hanya web responsive di awal).
- Integrasi ERP/akuntansi pihak ketiga.

---

## 3. Target Pengguna

| Persona | Deskripsi |
|---|---|
| **Customer** | Pengguna yang membeli produk atau memesan jasa/layanan |
| **Admin/Owner** | Pemilik toko, mengelola katalog, pesanan, konten |
| **Staff (opsional)** | Mengelola order & booking harian (role terbatas) |

---

## 4. Ruang Lingkup Fitur

### 4.1 Modul Customer (Frontend Publik)
- Registrasi/Login (email, opsional login sosial)
- Katalog produk (list, detail, filter, search, kategori)
- Katalog jasa (list, detail, jadwal/slot tersedia)
- Keranjang belanja (produk)
- Sistem booking (jasa) — pilih tanggal & slot waktu
- Checkout & pembayaran (Midtrans/Xendit — payment gateway lokal)
- Riwayat pesanan & status tracking
- Riwayat booking jasa & status (pending/confirmed/selesai/batal)
- Ulasan & rating (produk dan jasa)
- Wishlist (opsional v1.1)
- Notifikasi (email/在-app) status pesanan/booking
- Profil pengguna & alamat pengiriman

### 4.2 Modul Admin (Dashboard)
- Manajemen produk (CRUD, kategori, varian, stok, harga)
- Manajemen jasa (CRUD, durasi, kapasitas slot, jadwal kerja)
- Manajemen pesanan (update status, cetak invoice, resi pengiriman)
- Manajemen booking (approve/reject, kalender)
- Manajemen pengguna & role
- Manajemen banner/promo/diskon/voucher
- Laporan penjualan (produk vs jasa, grafik, export)
- Manajemen media/upload gambar (via MinIO)

### 4.3 Modul Sistem/Backend Umum
- Autentikasi & otorisasi (JWT + refresh token, RBAC)
- Payment gateway integration
- Notifikasi (email via SMTP, opsional WhatsApp API)
- File storage (MinIO S3) untuk gambar produk, avatar, invoice PDF
- Logging & audit trail
- Rate limiting & API security

---

## 5. Arsitektur Sistem

```
┌─────────────────────┐        ┌──────────────────────┐
│      FRONTEND        │  API   │       BACKEND         │
│  (NuxtJS / Laravel)  │◄──────►│ (Golang Fiber/NestJS) │
└─────────────────────┘  REST/  └──────────┬───────────┘
                          JSON              │
                                            │
                     ┌──────────────────────┼──────────────────────┐
                     ▼                      ▼                      ▼
              ┌─────────────┐       ┌──────────────┐      ┌────────────────┐
              │  PostgreSQL  │       │  MinIO (S3)   │      │ Redis (cache/  │
              │  (database)  │       │ object storage│      │ queue/session) │
              └─────────────┘       └──────────────┘      └────────────────┘
                                            │
                                            ▼
                                   ┌──────────────────┐
                                   │ Payment Gateway   │
                                   │ (Midtrans/Xendit) │
                                   └──────────────────┘
```

Komunikasi FE ↔ BE sepenuhnya via **REST API** (atau opsional GraphQL bila NestJS dipilih). BE tidak melakukan render halaman (headless), FE bertanggung jawab penuh atas UI/UX & SEO (SSR).

---

## 6. Tech Stack (Final)

| Layer | Pilihan Final | Catatan |
|---|---|---|
| **Backend** | **Golang + Fiber** | Prioritas performa & concurrency (penting untuk booking slot jasa & long-term scaling sebagai SaaS) |
| **Frontend** | **NuxtJS (Vue 3)** | SSR untuk SEO produk, konsisten dengan arsitektur headless |
| **Storage** | **MinIO (S3-compatible)** | Self-hosted, jalan sebagai container di VPS yang sama |
| **Payment Gateway** | **Midtrans** | Snap/Core API, mendukung banyak metode bayar lokal (VA, e-wallet, QRIS, kartu) |
| **Deployment** | **VPS + Docker Compose** | Semua service (BE, FE, DB, MinIO, Redis, Nginx) dijalankan sebagai container dalam satu `docker-compose.yml` |

### 6.1 Backend — Golang + Fiber
- Framework HTTP: **Fiber v2**
- ORM: **GORM** (dengan PostgreSQL driver)
- Validasi: `go-playground/validator`
- JWT: `golang-jwt/jwt`
- Storage SDK: `minio-go`
- Struktur project: modular (per domain: `product`, `service`, `order`, `booking`, `payment`, `user`), mengikuti pola **clean architecture ringan** (handler → usecase/service → repository) agar mudah dites dan direfactor ke multi-tenant nanti.

### 6.2 Frontend — NuxtJS
- Nuxt 3, Vue 3 Composition API
- State management: **Pinia**
- Styling: TailwindCSS
- SSR aktif untuk halaman katalog produk/jasa (SEO), CSR untuk dashboard admin (SPA mode)
- Integrasi Midtrans Snap.js di sisi client untuk popup pembayaran

### 6.3 Storage — MinIO
- Dijalankan sebagai container terpisah dalam Docker Compose, volume persistent di VPS
- Bucket: `products`, `services`, `avatars`, `invoices`, `banners`
- **Konvensi path object: `{tenant_slug}/{bucket}/{filename}`** — meskipun v1 hanya 1 tenant (misal `default`), konvensi ini disiapkan agar migrasi ke multi-tenant tidak perlu memindahkan/rename file lama (lihat Bab 13)
- Akses via **presigned URL**, upload langsung dari FE ke MinIO (tidak lewat BE) untuk mengurangi beban server

### 6.4 Payment Gateway — Midtrans
- Gunakan **Snap API** (redirect/popup) untuk kecepatan integrasi di v1
- Webhook `POST /api/payments/webhook/midtrans` untuk notifikasi status transaksi (`settlement`, `pending`, `expire`, `deny`, `cancel`)
- Order harus idempotent terhadap webhook yang terkirim berulang (Midtrans bisa retry)
- Simpan `transaction_id`, `order_id`, `gross_amount`, `payment_type`, `status`, `raw_response` (JSON) di tabel `payments` untuk audit

### 6.5 Deployment — VPS + Docker Compose
Contoh susunan service dalam `docker-compose.yml`:

```yaml
services:
  backend:      # Golang Fiber app
  frontend:     # NuxtJS app (Node runtime / atau static+nginx bila SSG)
  postgres:     # Database
  redis:        # Cache & queue
  minio:        # Object storage
  nginx:        # Reverse proxy + SSL (Let's Encrypt/Certbot)
```
- **Nginx** sebagai reverse proxy di depan FE & BE, sekaligus terminasi SSL
- Environment variables dikelola via `.env` per service (jangan commit ke repo)
- Backup otomatis: `pg_dump` terjadwal (cron) + `mc mirror` untuk backup bucket MinIO ke storage eksternal/VPS lain
- Deployment manual/simple CI: `git pull` → `docker compose up -d --build` (bisa ditingkatkan ke GitHub Actions self-hosted runner nanti)

### 6.6 Stack Pendukung Lainnya
| Komponen | Pilihan |
|---|---|
| Database | PostgreSQL |
| Cache/Queue | Redis (session, cache katalog, queue notifikasi/email) |
| CI/CD | Manual di awal → GitHub Actions (opsional, saat siap) |
| Monitoring | Sentry (error tracking), opsional Prometheus + Grafana saat traffic naik |

---

## 7. Data Model (High Level)

> Catatan: setiap tabel di bawah menyertakan kolom `tenant_id` sejak v1 (lihat Bab 14 — Strategi SaaS-Readiness), meskipun nilainya selalu tenant default untuk saat ini.

- **Tenant** (id, slug, name, status) — v1 hanya berisi 1 baris (`default`)
- **User** (id, tenant_id, name, email, password_hash, role, phone)
- **Product** (id, tenant_id, name, description, price, stock, sku, category_id, images[])
- **ProductVariant** (id, product_id, name, price_diff, stock)
- **Service** (id, tenant_id, name, description, price, duration_minutes, category_id, images[])
- **ServiceSlot** (id, service_id, date, start_time, end_time, capacity, booked_count)
- **Category** (id, tenant_id, name, type: product|service)
- **Order** (id, tenant_id, user_id, type: product|service, status, total, payment_status)
- **OrderItem** (id, order_id, product_id/service_id, qty, price)
- **Booking** (id, order_id, service_slot_id, status: pending|confirmed|done|cancelled)
- **Payment** (id, order_id, gateway, status, transaction_id, gross_amount, payment_type, raw_response, paid_at)
- **Review** (id, user_id, product_id/service_id, rating, comment)
- **Voucher** (id, tenant_id, code, discount_type, value, expiry)
- **MediaFile** (id, tenant_id, bucket, object_key, url, uploaded_by)
- **Setting** (id, tenant_id, key, value) — menyimpan konfigurasi toko (nama, logo, kontak) agar tidak hardcode di kode

---

## 8. Alur Utama (User Flow)

### 8.1 Alur Pembelian Produk (v1 — berhenti di pembayaran)
1. Customer browse katalog → tambah ke keranjang
2. Checkout → isi alamat → pilih metode bayar
3. Redirect/popup **Midtrans Snap** → bayar
4. Webhook Midtrans → update status order jadi `paid`
5. Admin memproses pengiriman **secara manual di luar sistem** (belum ada modul delivery di v1)
6. Customer terima notifikasi order `paid` & bisa beri ulasan

> Catatan: Update resi/tracking pengiriman **tidak** termasuk v1. Order dianggap "selesai dari sisi sistem" setelah status `paid`. Modul pengiriman & armada akan menyusul di Fase 5 (lihat Bab 11).

### 8.2 Alur Booking Jasa
1. Customer pilih jasa → lihat kalender slot tersedia
2. Pilih tanggal & jam → isi detail kebutuhan
3. Checkout & bayar (atau DP)
4. Admin/staff approve booking
5. Notifikasi konfirmasi ke customer (H-1 reminder opsional)
6. Setelah selesai → status "done", customer bisa review

---

## 9. Non-Functional Requirements

| Kategori | Requirement |
|---|---|
| **Performa** | API response < 300ms untuk 95% request (non-payment) |
| **Skalabilitas** | Backend stateless, bisa horizontal scaling |
| **Keamanan** | HTTPS wajib, JWT dengan refresh token, hashing bcrypt/argon2, rate limiting, validasi input ketat, presigned URL utk file |
| **Ketersediaan** | Target uptime 99% (untuk personal project) |
| **Kompatibilitas** | Responsive (mobile-first), mendukung browser modern |
| **Localization** | Bahasa Indonesia sebagai default, mata uang IDR |
| **Backup** | Backup database harian, backup MinIO bucket berkala |

---

## 10. API Design (Contoh Endpoint Utama)

```
Auth
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/refresh

Produk
GET    /api/products
GET    /api/products/:id
POST   /api/admin/products
PUT    /api/admin/products/:id
DELETE /api/admin/products/:id

Jasa
GET    /api/services
GET    /api/services/:id/slots
POST   /api/admin/services
POST   /api/admin/services/:id/slots

Order & Booking
POST   /api/orders               (checkout produk)
POST   /api/bookings             (checkout jasa)
GET    /api/orders/:id
PATCH  /api/admin/orders/:id/status

Upload (MinIO)
POST   /api/uploads/presign      (generate presigned URL utk upload langsung ke MinIO)

Payment
POST   /api/payments/webhook/midtrans   (callback notifikasi status transaksi dari Midtrans)

Settings
GET    /api/settings             (branding/konfigurasi toko, dipakai FE untuk tampilan)
PUT    /api/admin/settings        (update konfigurasi toko)
```

---

## 11. Roadmap / Milestone

| Fase | Cakupan | Estimasi |
|---|---|---|
| **Fase 1 – MVP** | Auth, katalog produk, keranjang, checkout, integrasi Midtrans (sampai status `paid`), admin dasar | 4–6 minggu |
| **Fase 2** | Modul jasa/booking + kalender slot | 3–4 minggu |
| **Fase 3** | Voucher, review, laporan penjualan | 2–3 minggu |
| **Fase 4** | Optimasi performa, SEO, monitoring, hardening security | 2 minggu |
| **Fase 5 – Delivery & Fleet Management** *(lanjutan, di luar v1)* | Modul pengiriman order + manajemen armada (lihat detail di bawah) | 4–6 minggu |
| **Fase 6 – SaaS Activation** *(opsional, jika dibutuhkan)* | Aktivasi multi-tenant: onboarding tenant baru, subdomain, billing SaaS | 4–8 minggu |

### 11.1 Detail Fase 5 — Delivery & Fleet Management (Rencana Lanjutan)

Modul ini menyusul setelah v1 (yang berhenti di payment) berjalan stabil. Cakupan yang direncanakan:

**Delivery Management**
- Status pengiriman: `waiting_pickup` → `on_the_way` → `delivered` / `failed`
- Assign order ke kurir/driver
- Input bukti pengiriman (foto/tanda tangan) — disimpan di MinIO bucket `pod` (proof of delivery)
- Notifikasi status ke customer (email/WA)
- Opsional: perhitungan ongkir otomatis (integrasi RajaOngkir/Biteship) atau ongkir internal (armada sendiri)

**Fleet / Manajemen Armada**
- Data kendaraan (jenis, plat nomor, kapasitas, status: aktif/maintenance)
- Data driver (nama, kontak, kendaraan yang di-assign, status: available/on-duty/off)
- Penjadwalan/assignment order ke driver+kendaraan (manual dulu, bisa otomatis nanti berdasar rute/kapasitas)
- Riwayat perjalanan per kendaraan/driver
- Opsional (v2+): tracking lokasi real-time via GPS (mobile app driver terpisah)

**Entitas tambahan yang akan muncul di Fase 5:**
- `Vehicle` (id, plate_number, type, capacity, status)
- `Driver` (id, user_id, vehicle_id, status)
- `Delivery` (id, order_id, driver_id, vehicle_id, status, pod_image_url, delivered_at)

> Karena ini fase lanjutan, detail requirement lengkap (termasuk apakah butuh mobile app driver) akan dibuatkan PRD/dokumen teknis terpisah saat Fase 5 dimulai.

---

## 12. Metrik Keberhasilan (Success Metrics)

- Waktu checkout rata-rata < 2 menit
- Tingkat konversi visitor → order minimal 2%
- Tidak ada downtime pembayaran (payment webhook reliability > 99%)
- Waktu approve booking oleh admin < 24 jam

---

## 13. Risiko & Mitigasi

| Risiko | Mitigasi |
|---|---|
| Double booking slot jasa | Lock/transaction saat booking + validasi capacity |
| Payment gateway callback gagal | Retry mechanism + reconciliation job berkala |
| File upload besar membebani server | Upload langsung ke MinIO via presigned URL (bypass BE) |
| Stok race condition saat checkout bersamaan | Gunakan DB transaction + row locking |

---

## 14. Strategi SaaS-Readiness (Single-Tenant → Multi-Tenant di Masa Depan)

Karena v1 dijalankan single-tenant tapi punya ambisi jadi SaaS jangka panjang, berikut aturan desain yang diterapkan **sejak awal** agar migrasi nanti murah (tidak perlu rewrite):

### 13.1 Database
- Setiap tabel utama (`products`, `services`, `orders`, `bookings`, `users`, dll.) **sudah punya kolom `tenant_id`** sejak v1, meski nilainya selalu diisi 1 (tenant default).
- Semua query di repository layer **wajib** difilter berdasarkan `tenant_id` (walau saat ini hasilnya selalu sama), bukan diasumsikan implicit "cuma ada 1 toko". Ini memastikan saat multi-tenant diaktifkan, tidak ada query yang "bocor" data antar tenant.
- Index composite `(tenant_id, ...)` disiapkan dari awal untuk performa saat data tenant bertambah.

### 13.2 Autentikasi & Otorisasi
- JWT payload menyertakan `tenant_id` sejak awal (bukan ditambah belakangan), walau nilainya statis untuk v1.
- Middleware auth di Fiber sudah dirancang untuk "resolve tenant" dari JWT/domain — sehingga saat multi-tenant aktif, tinggal mengubah sumber resolusi tenant (dari statis ke dinamis berdasarkan subdomain/header).

### 13.3 Struktur Kode (Backend)
- Business logic dipisah per domain module (bukan satu file besar), sehingga penambahan lapisan tenant-scoping tidak menyentuh seluruh codebase, cukup di layer repository/middleware.
- Hindari hardcode asumsi "hanya ada 1 toko" di layer service (misal nama toko, logo, kontak) — simpan sebagai data di tabel `settings` (per tenant), bukan konstanta di kode.

### 13.4 Storage (MinIO)
- Konvensi path `{tenant_slug}/{bucket}/{filename}` sudah dipakai sejak v1 (tenant_slug = `default` untuk saat ini).
- Ini membuat migrasi ke multi-tenant tidak memerlukan pemindahan/rename file lama — tenant baru otomatis mendapat "folder" terpisah.

### 13.5 Frontend (Nuxt)
- Konfigurasi tema/branding (logo, warna, nama toko) diambil dari API `/api/settings`, bukan di-hardcode di komponen — sehingga saat multi-tenant aktif, FE otomatis bisa menampilkan branding berbeda per tenant tanpa perubahan kode, cukup beda data.

### 13.6 Yang SENGAJA DITUNDA (agar tidak over-engineering di v1)
- Belum ada UI onboarding tenant baru / pendaftaran seller.
- Belum ada sistem billing/subscription SaaS (baru relevan saat Fase 6 diaktifkan).
- Belum ada routing subdomain per tenant (`{tenant}.domain.com`) — baru diimplementasikan saat Fase 6.
- Belum perlu database terpisah per tenant (schema-per-tenant/DB-per-tenant) — cukup shared-schema dengan `tenant_id` untuk skala saat ini.

> Prinsipnya: **desain data & kode "tenant-aware" sejak hari pertama**, tapi **fitur multi-tenant-nya sendiri baru diaktifkan saat benar-benar dibutuhkan** (Fase 6). Ini menghindari kompleksitas premature tapi tetap menghindari migrasi mahal di kemudian hari.

---

*Dokumen ini adalah draft awal dan dapat disesuaikan sesuai kebutuhan lebih lanjut (misalnya penambahan multi-bahasa, program afiliasi, atau integrasi ekspedisi pengiriman).*
## 15. Keputusan yang Sudah Diambil (Final)

| Keputusan | Pilihan |
|---|---|
| Backend | **Golang + Fiber** |
| Frontend | **NuxtJS** |
| Payment gateway | **Midtrans** |
| Deployment | **VPS + Docker Compose** |
| Cakupan v1 | Sampai status **paid**, tanpa modul delivery |
| Arah SaaS | **Single-tenant sekarang**, arsitektur disiapkan agar mudah direfactor jadi multi-tenant (lihat Bab 13) |

---

