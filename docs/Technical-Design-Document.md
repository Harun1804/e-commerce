# Technical Design Document (TDD)
# E-Commerce Platform — Produk & Jasa

| | |
|---|---|
| **Versi** | 1.0 |
| **Tanggal** | 7 Juli 2026 |
| **Terkait** | PRD-Ecommerce-Personal.md v1.0 |
| **Stack** | Golang Fiber, NuxtJS, PostgreSQL, MinIO, Redis, Midtrans |

Dokumen ini adalah turunan teknis dari PRD, berisi detail implementasi: skema database, struktur project, konfigurasi deployment, dan setup replikasi.

---

## 1. Skema Database (PostgreSQL)

### 1.1 Prinsip Desain
- Semua tabel domain punya `tenant_id` (lihat PRD Bab 14 — SaaS-Readiness), default mengarah ke 1 baris di tabel `tenants`.
- Primary key menggunakan `UUID` (bukan auto-increment int) — lebih aman untuk exposed ke API publik dan lebih siap untuk skenario multi-tenant/distributed di masa depan.
- Soft-delete via kolom `deleted_at` (nullable timestamp) di tabel-tabel penting, konsisten dengan konvensi GORM (`gorm.DeletedAt`).
- Exclusion constraint dipakai di `service_bookings` untuk mencegah **double-booking** langsung di level database (bukan hanya application-level lock).

### 1.2 DDL Lengkap

```sql
-- Extension yang dibutuhkan
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "btree_gist"; -- untuk exclusion constraint pada range waktu

-- =========================
-- TENANT & SETTINGS
-- =========================
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug VARCHAR(100) UNIQUE NOT NULL DEFAULT 'default',
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active | suspended
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    key VARCHAR(100) NOT NULL,
    value JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, key)
);

-- =========================
-- USERS & AUTH
-- =========================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(30),
    role VARCHAR(20) NOT NULL DEFAULT 'customer', -- customer | admin | staff
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (tenant_id, email)
);
CREATE INDEX idx_users_tenant ON users(tenant_id);

CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    label VARCHAR(50),
    recipient_name VARCHAR(255) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    full_address TEXT NOT NULL,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =========================
-- CATALOG: PRODUCT
-- =========================
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(150) NOT NULL,
    type VARCHAR(20) NOT NULL, -- product | service
    slug VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, slug)
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    category_id UUID REFERENCES categories(id),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(14,2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    sku VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active | draft | archived
    images JSONB NOT NULL DEFAULT '[]', -- array of MinIO object keys
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (tenant_id, slug)
);
CREATE INDEX idx_products_tenant_status ON products(tenant_id, status);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id),
    name VARCHAR(150) NOT NULL, -- misal "Ukuran L / Merah"
    price_diff NUMERIC(14,2) NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    sku VARCHAR(100)
);

-- =========================
-- CATALOG: SERVICE
-- =========================
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    category_id UUID REFERENCES categories(id),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(14,2) NOT NULL,
    duration_minutes INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    images JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (tenant_id, slug)
);

CREATE TABLE service_slots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL REFERENCES services(id),
    slot_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 1,
    booked_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_service_slots_service_date ON service_slots(service_id, slot_date);

-- =========================
-- ORDERS (PRODUK)
-- =========================
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    user_id UUID NOT NULL REFERENCES users(id),
    order_number VARCHAR(30) UNIQUE NOT NULL,
    order_type VARCHAR(20) NOT NULL, -- product | service
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    -- pending -> awaiting_payment -> paid -> (v1 berhenti di sini) / expired / cancelled
    subtotal NUMERIC(14,2) NOT NULL,
    discount NUMERIC(14,2) NOT NULL DEFAULT 0,
    total NUMERIC(14,2) NOT NULL,
    shipping_address_id UUID REFERENCES addresses(id),
    voucher_id UUID,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_orders_tenant_status ON orders(tenant_id, status);
CREATE INDEX idx_orders_user ON orders(user_id);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id),
    product_id UUID REFERENCES products(id),
    product_variant_id UUID REFERENCES product_variants(id),
    service_id UUID REFERENCES services(id),
    name_snapshot VARCHAR(255) NOT NULL, -- simpan nama saat transaksi (jaga2 produk berubah nama)
    price_snapshot NUMERIC(14,2) NOT NULL,
    qty INTEGER NOT NULL DEFAULT 1,
    subtotal NUMERIC(14,2) NOT NULL
);

-- =========================
-- BOOKING (JASA)
-- =========================
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id),
    service_slot_id UUID NOT NULL REFERENCES service_slots(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    -- pending -> confirmed -> done / cancelled
    customer_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =========================
-- PAYMENT (MIDTRANS)
-- =========================
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id),
    gateway VARCHAR(30) NOT NULL DEFAULT 'midtrans',
    midtrans_order_id VARCHAR(100) NOT NULL,
    transaction_id VARCHAR(100),
    payment_type VARCHAR(50), -- va, gopay, qris, credit_card, dll
    gross_amount NUMERIC(14,2) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    -- pending -> settlement/capture -> (paid) / expire / deny / cancel
    raw_response JSONB,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (midtrans_order_id)
);
CREATE INDEX idx_payments_order ON payments(order_id);

-- =========================
-- REVIEW & VOUCHER
-- =========================
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    product_id UUID REFERENCES products(id),
    service_id UUID REFERENCES services(id),
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE vouchers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    code VARCHAR(50) NOT NULL,
    discount_type VARCHAR(20) NOT NULL, -- percentage | fixed
    value NUMERIC(14,2) NOT NULL,
    max_usage INTEGER,
    used_count INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, code)
);

-- =========================
-- MEDIA (MinIO reference)
-- =========================
CREATE TABLE media_files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    bucket VARCHAR(50) NOT NULL,
    object_key TEXT NOT NULL, -- format: {tenant_slug}/{bucket}/{filename}
    url TEXT NOT NULL,
    uploaded_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 1.3 Mencegah Double-Booking di Level Database

Selain validasi di application layer, tambahkan **exclusion constraint** agar dua booking tidak bisa saling tumpang tindih pada slot yang sama meski terjadi race condition:

```sql
-- Tambahan kolom range waktu di service_slots (generated, untuk exclusion constraint)
ALTER TABLE service_slots
    ADD COLUMN time_range tsrange
    GENERATED ALWAYS AS (
        tsrange(
            (slot_date + start_time)::timestamp,
            (slot_date + end_time)::timestamp
        )
    ) STORED;

-- Cegah slot dengan waktu tumpang tindih untuk service yang sama
ALTER TABLE service_slots
    ADD CONSTRAINT no_overlapping_slots
    EXCLUDE USING gist (
        service_id WITH =,
        time_range WITH &&
    );
```

> Untuk mencegah **over-booking dalam 1 slot yang sama** (misal capacity 3 tapi ada 4 booking masuk bersamaan), gunakan `SELECT ... FOR UPDATE` pada baris `service_slots` saat proses booking di dalam transaction Golang (lihat Bab 3.2).

---

## 2. Struktur Project Backend (Golang + Fiber)

Mengikuti **clean architecture ringan**: `handler → usecase → repository`, dipisah per domain module.

```
backend/
├── cmd/
│   └── api/
│       └── main.go                  # entrypoint, load config, wire dependencies
├── internal/
│   ├── config/
│   │   └── config.go                 # load env vars (viper/envconfig)
│   ├── middleware/
│   │   ├── auth.go                   # JWT verify + resolve tenant_id dari token
│   │   ├── tenant_scope.go           # inject tenant_id ke context, dipakai semua repo
│   │   ├── rate_limit.go
│   │   └── logger.go
│   ├── domain/
│   │   ├── product/
│   │   │   ├── entity.go
│   │   │   ├── handler.go
│   │   │   ├── usecase.go
│   │   │   └── repository.go
│   │   ├── service/                  # domain "jasa", beda dgn folder service teknis
│   │   │   ├── entity.go
│   │   │   ├── handler.go
│   │   │   ├── usecase.go
│   │   │   └── repository.go
│   │   ├── order/
│   │   ├── booking/
│   │   ├── payment/
│   │   │   ├── entity.go
│   │   │   ├── handler.go
│   │   │   ├── usecase.go
│   │   │   ├── repository.go
│   │   │   └── midtrans_client.go    # wrapper Snap/Core API + webhook parser
│   │   ├── user/
│   │   └── settings/
│   ├── pkg/
│   │   ├── storage/
│   │   │   └── minio_client.go       # presigned URL generator, upload/delete
│   │   ├── database/
│   │   │   └── postgres.go           # gorm.Open + connection pool config
│   │   ├── cache/
│   │   │   └── redis_client.go
│   │   └── validator/
│   │       └── validator.go
│   └── router/
│       └── router.go                 # daftar semua route Fiber, grup per domain
├── migrations/                       # SQL migration files (golang-migrate)
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── Dockerfile
├── go.mod
└── go.sum
```

**Catatan desain penting:**
- `middleware/tenant_scope.go` meng-inject `tenant_id` ke `context.Context` dari awal request. Semua fungsi di `repository.go` **wajib** menerima `ctx` dan filter query berdasarkan `tenant_id` dari context — bukan hardcode/asumsi 1 tenant. Ini persis mekanisme yang disebut di PRD Bab 14.
- `payment/midtrans_client.go` menangani verifikasi signature webhook Midtrans (`signature_key` = SHA512 dari `order_id+status_code+gross_amount+ServerKey`) sebelum memproses notifikasi — wajib untuk keamanan webhook.

---

## 3. Struktur Project Frontend (NuxtJS)

```
frontend/
├── app/
├── assets/
│   └── css/
├── components/
│   ├── product/
│   ├── service/
│   ├── cart/
│   ├── booking/
│   └── ui/                           # button, modal, dsb (design system dasar)
├── composables/
│   ├── useAuth.ts
│   ├── useCart.ts
│   ├── useCheckout.ts
│   └── useSettings.ts                # fetch branding dari /api/settings
├── layouts/
│   ├── default.vue                   # layout customer-facing
│   └── admin.vue                     # layout dashboard admin
├── middleware/
│   └── auth.global.ts
├── pages/
│   ├── index.vue
│   ├── products/
│   │   ├── index.vue
│   │   └── [slug].vue
│   ├── services/
│   │   ├── index.vue
│   │   └── [slug].vue
│   ├── cart.vue
│   ├── checkout.vue
│   └── admin/
│       ├── products/
│       ├── services/
│       └── orders/
├── stores/                           # Pinia
│   ├── cart.ts
│   ├── auth.ts
│   └── settings.ts
├── nuxt.config.ts
├── Dockerfile
└── package.json
```

### 3.1 Integrasi Midtrans Snap di Nuxt
```html
<!-- nuxt.config.ts: daftarkan script Snap.js -->
app: {
  head: {
    script: [
      { src: 'https://app.sandbox.midtrans.com/snap/snap.js',
        'data-client-key': process.env.MIDTRANS_CLIENT_KEY }
    ]
  }
}
```
Alur: BE membuat transaksi via Core/Snap API → BE kirim `snap_token` ke FE → FE panggil `window.snap.pay(snap_token)` → popup pembayaran muncul → setelah selesai, FE polling/redirect ke halaman status order (status final tetap ditentukan webhook, bukan response client-side, karena client-side bisa dimanipulasi).

---

## 4. Docker Compose (VPS Deployment)

```yaml
version: "3.9"

services:
  postgres-primary:
    image: postgres:16
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - pg_primary_data:/var/lib/postgresql/data
      - ./postgres/primary/postgresql.conf:/etc/postgresql/postgresql.conf
      - ./postgres/primary/pg_hba.conf:/etc/postgresql/pg_hba.conf
    command: postgres -c config_file=/etc/postgresql/postgresql.conf
    ports:
      - "5432:5432"
    networks:
      - internal

  postgres-replica:
    image: postgres:16
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - pg_replica_data:/var/lib/postgresql/data
    depends_on:
      - postgres-primary
    entrypoint: ["/bin/bash", "/scripts/init-replica.sh"]
    volumes:
      - ./postgres/replica/init-replica.sh:/scripts/init-replica.sh
      - pg_replica_data:/var/lib/postgresql/data
    networks:
      - internal

  redis:
    image: redis:7-alpine
    restart: always
    volumes:
      - redis_data:/data
    networks:
      - internal

  minio:
    image: minio/minio
    restart: always
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ROOT_USER}
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD}
    volumes:
      - minio_data:/data
    ports:
      - "9000:9000"
      - "9001:9001"
    networks:
      - internal

  backend:
    build: ./backend
    restart: always
    env_file: ./backend/.env
    depends_on:
      - postgres-primary
      - redis
      - minio
    networks:
      - internal

  frontend:
    build: ./frontend
    restart: always
    env_file: ./frontend/.env
    depends_on:
      - backend
    networks:
      - internal

  nginx:
    image: nginx:alpine
    restart: always
    volumes:
      - ./nginx/conf.d:/etc/nginx/conf.d
      - ./nginx/certbot/conf:/etc/letsencrypt
      - ./nginx/certbot/www:/var/www/certbot
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - backend
      - frontend
    networks:
      - internal

  certbot:
    image: certbot/certbot
    volumes:
      - ./nginx/certbot/conf:/etc/letsencrypt
      - ./nginx/certbot/www:/var/www/certbot
    entrypoint: "/bin/sh -c 'trap exit TERM; while :; do certbot renew; sleep 12h & wait $${!}; done;'"

volumes:
  pg_primary_data:
  pg_replica_data:
  redis_data:
  minio_data:

networks:
  internal:
    driver: bridge
```

> Catatan: `backend` melakukan **write** ke `postgres-primary`. Untuk skala lanjut, tambahkan koneksi kedua di config Golang untuk **read** dari `postgres-replica` pada endpoint read-heavy (katalog produk/jasa), dengan fallback ke primary bila replica down.

---

## 5. Setup Replikasi PostgreSQL (Primary → Replica)

### 5.1 Primary — `postgresql.conf`
```conf
listen_addresses = '*'
wal_level = replica
max_wal_senders = 5
wal_keep_size = 512MB
hot_standby = on
```

### 5.2 Primary — `pg_hba.conf`
```conf
# Izinkan replica melakukan streaming replication
host    replication     ${DB_USER}      0.0.0.0/0      md5
```

### 5.3 Script Inisialisasi Replica — `init-replica.sh`
```bash
#!/bin/bash
set -e

if [ -z "$(ls -A /var/lib/postgresql/data)" ]; then
  echo "Cloning data dari primary..."
  pg_basebackup -h postgres-primary -D /var/lib/postgresql/data \
    -U ${DB_USER} -Fp -Xs -P -R
fi

exec docker-entrypoint.sh postgres
```

Flag `-R` otomatis membuat `standby.signal` dan mengisi `primary_conninfo` di `postgresql.auto.conf`, sehingga replica langsung tahu harus streaming dari `postgres-primary`.

### 5.4 Verifikasi Replikasi Berjalan
```sql
-- Di primary:
SELECT client_addr, state, sync_state FROM pg_stat_replication;

-- Di replica:
SELECT pg_is_in_recovery(); -- harus return true
```

### 5.5 Catatan Failover
Setup di atas adalah **replikasi read-only sederhana** (belum auto-failover). Untuk auto-failover production-grade nanti (saat traffic naik), pertimbangkan menambahkan **Patroni + etcd** atau **repmgr** sebagai peningkatan — tapi untuk kebutuhan v1 (personal project, traffic belum tinggi), setup manual streaming replication ini sudah cukup memberi **read replica untuk scaling baca** dan **backup hangat** jika primary down (butuh promote manual: `pg_ctl promote`).

---

## 6. Environment Variables (Ringkasan)

```env
# Backend (.env)
APP_ENV=production
APP_PORT=8080
JWT_SECRET=
JWT_EXPIRE_MINUTES=60

DB_HOST=postgres-primary
DB_HOST_REPLICA=postgres-replica
DB_PORT=5432
DB_USER=
DB_PASSWORD=
DB_NAME=ecommerce

REDIS_HOST=redis
REDIS_PORT=6379

MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=
MINIO_SECRET_KEY=
MINIO_USE_SSL=false

MIDTRANS_SERVER_KEY=
MIDTRANS_CLIENT_KEY=
MIDTRANS_IS_PRODUCTION=false

# Frontend (.env)
NUXT_PUBLIC_API_BASE=https://api.domainkamu.com
NUXT_PUBLIC_MIDTRANS_CLIENT_KEY=
```

---

## 7. Alur Teknis Kritis

### 7.1 Booking Slot (Anti Race-Condition)
```sql
BEGIN;
SELECT * FROM service_slots WHERE id = $1 FOR UPDATE;
-- cek booked_count < capacity di application code
UPDATE service_slots SET booked_count = booked_count + 1 WHERE id = $1;
INSERT INTO bookings (...) VALUES (...);
COMMIT;
```

### 7.2 Webhook Midtrans (Idempotent)
1. Terima payload webhook di `POST /api/payments/webhook/midtrans`.
2. Verifikasi `signature_key` (SHA512).
3. Cari `payments` row via `midtrans_order_id` — **jika status sudah `paid`, abaikan** (idempotent, karena Midtrans bisa kirim notifikasi berkali-kali).
4. Update `payments.status` dan `orders.status` dalam satu DB transaction.
5. Return HTTP 200 ke Midtrans secepatnya (proses notifikasi lanjutan seperti kirim email bisa lewat queue/goroutine async, jangan blocking response).

### 7.3 Upload File (Presigned URL)
1. FE minta presigned URL: `POST /api/uploads/presign` dengan `{bucket, filename, content_type}`.
2. BE generate presigned PUT URL dari MinIO (`minio-go`), scoped ke path `{tenant_slug}/{bucket}/{filename}`.
3. FE upload langsung ke MinIO pakai URL tsb (tidak lewat BE).
4. FE kirim konfirmasi ke BE (`POST /api/media`) untuk menyimpan record di tabel `media_files`.

---

## 8. Checklist Sebelum Go-Live (v1)

- [ ] Migration schema sudah dijalankan (`golang-migrate up`)
- [ ] `pg_hba.conf` primary sudah restrict IP (jangan `0.0.0.0/0` di production, ganti ke IP replica spesifik)
- [ ] Replikasi primary→replica sudah diverifikasi (`pg_stat_replication`)
- [ ] Webhook Midtrans sudah diuji di sandbox, termasuk skenario retry/duplicate notification
- [ ] Bucket MinIO sudah dibuat + policy akses sudah diset (public-read untuk gambar produk, private untuk invoice/POD)
- [ ] SSL certificate (certbot) aktif di Nginx
- [ ] Backup terjadwal (`pg_dump` + `mc mirror`) sudah di-cron
- [ ] Rate limiting aktif di endpoint auth & checkout
