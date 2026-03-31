# WE SAVING API

Backend API sederhana untuk aplikasi pencatatan keuangan pribadi: autentikasi, kategori pengeluaran, alokasi budget per kategori, dan pencatatan salary. Project ini ditulis dengan Go, Chi Router, PostgreSQL, dan JWT.

## Fitur yang Sudah Tersedia

- Register, login, dan endpoint `me`
- CRUD kategori pengeluaran
- CRUD budget per kategori
- CRUD salary + cek total salary
- Validasi request dan response JSON yang konsisten
- Migration database berbasis SQL
- Unit test dan sebagian integration test

## Status Fitur

Endpoint yang sudah aktif saat ini:

- Auth
- Categories
- Category Budgets
- Salaries

Schema database yang sudah ada tetapi endpoint-nya belum diregistrasi:

- Savings
- Transactions

## Tech Stack

- Go `1.25`
- PostgreSQL `15`
- [`chi`](https://github.com/go-chi/chi)
- [`sqlx`](https://github.com/jmoiron/sqlx)
- [`golang-jwt/jwt`](https://github.com/golang-jwt/jwt)
- [`validator`](https://github.com/go-playground/validator)
- [`golang-migrate`](https://github.com/golang-migrate/migrate)
- Docker + Docker Compose

## Struktur Folder

```text
.
├── cmd/we-saving-api           # entrypoint aplikasi
├── internal/app                # router, middleware, server bootstrap
├── internal/domain             # handler, service, repository per domain
├── internal/Infrastructures    # config dan koneksi database
├── internal/shared             # helper JWT, response, validation
├── migrations                  # SQL migration
├── Dockerfile
└── docker-compose.yml
```

## Prasyarat

Salah satu dari dua mode berikut:

1. Lokal tanpa Docker
   - Go `1.25+`
   - PostgreSQL `15+`
   - CLI `migrate` dari golang-migrate
2. Docker
   - Docker
   - Docker Compose

## Environment Variables

Copy template dulu:

```bash
cp .env.example .env
```

Variabel yang dipakai:

| Variable | Wajib | Keterangan |
| --- | --- | --- |
| `PORT` | Tidak | Port aplikasi. Default `8080` |
| `JWT_SECRET` | Ya | Secret untuk sign JWT |
| `DATABASE_URL` | Ya | Connection string PostgreSQL untuk jalan lokal |
| `DB_USER` | Ya jika pakai Docker Compose | User PostgreSQL container |
| `DB_PASSWORD` | Ya jika pakai Docker Compose | Password PostgreSQL container |
| `DB_NAME` | Ya jika pakai Docker Compose | Nama database PostgreSQL container |
| `TEST_DATABASE_URL` | Tidak | Dipakai untuk integration test |

Catatan:

- Untuk jalan lokal, `DATABASE_URL` biasanya memakai host `localhost`.
- Untuk Docker Compose, file `docker-compose.yml` sudah override koneksi app ke host `db`, jadi `DATABASE_URL` di `.env` boleh tetap versi lokal.

## Menjalankan Secara Lokal

1. Copy env:

```bash
cp .env.example .env
```

2. Buat database sesuai `.env`, misalnya:

```bash
createdb we_saving
createdb we_saving_test
```

3. Jalankan migration:

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

4. Jalankan aplikasi:

```bash
go run ./cmd/we-saving-api
```

Server akan jalan di:

```text
http://localhost:8080
```

Base path API:

```text
http://localhost:8080/api/v1
```

## Menjalankan Dengan Docker Compose

1. Copy env:

```bash
cp .env.example .env
```

2. Jalankan container:

```bash
docker compose up --build
```

App akan expose port `8080`, dan PostgreSQL expose port `5432`.

Untuk menghentikan:

```bash
docker compose down
```

Kalau mau sekalian hapus volume database:

```bash
docker compose down -v
```

## Menjalankan Test

Unit test:

```bash
go test ./...
```

Integration test membutuhkan `TEST_DATABASE_URL` mengarah ke database test yang bisa diakses:

```bash
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/we_saving_test?sslmode=disable" go test -tags=integration ./...
```

Kalau environment kamu membatasi cache default Go, gunakan cache lokal project:

```bash
GOCACHE=$(pwd)/.cache/go-build go test ./...
```

## Format Response

Mayoritas endpoint mengembalikan bentuk JSON seperti ini:

```json
{
  "code": 200,
  "status": "OK",
  "message": "success",
  "data": {},
  "error": false
}
```

Response validasi biasanya menambahkan field `valid`:

```json
{
  "code": 400,
  "status": "NOK",
  "message": "Validation Failed",
  "error": true,
  "valid": {
    "username": "username is required"
  }
}
```

## Autentikasi

Semua endpoint selain register, login, dan `GET /api/v1/test` membutuhkan header:

```text
Authorization: Bearer <token>
```

## Ringkasan Endpoint

Base URL: `http://localhost:8080/api/v1`

### Public

| Method | Endpoint | Keterangan |
| --- | --- | --- |
| `GET` | `/test` | Endpoint cek server |
| `POST` | `/auth/register` | Register user baru |
| `POST` | `/auth/login` | Login dan ambil JWT |

### Protected

| Method | Endpoint | Keterangan |
| --- | --- | --- |
| `GET` | `/auth/me` | Ambil profil user login |
| `POST` | `/categories/create` | Buat kategori |
| `GET` | `/categories/all` | Ambil semua kategori milik user |
| `GET` | `/categories/{id}` | Ambil detail kategori |
| `PUT` | `/categories/{id}` | Update kategori |
| `DELETE` | `/categories/{id}` | Hapus kategori |
| `POST` | `/category-budgets/create` | Buat budget kategori |
| `GET` | `/category-budgets/category/{id}` | Budget aktif per kategori |
| `GET` | `/category-budgets/category/{id}/all` | Semua riwayat budget per kategori |
| `PUT` | `/category-budgets/budget/{id}` | Update budget |
| `DELETE` | `/category-budgets/budget/{id}` | Hapus budget |
| `POST` | `/salary/create` | Buat salary |
| `GET` | `/salary/all` | Ambil semua salary user |
| `GET` | `/salary/check` | Cek status salary |
| `GET` | `/salary/total` | Total salary user |
| `GET` | `/salary/{id}` | Detail salary |
| `PUT` | `/salary/{id}` | Update salary |
| `DELETE` | `/salary/{id}` | Hapus salary |

## Contoh Request Cepat

Register:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "fahri",
    "fullname": "Fahriza",
    "email": "fahriza@example.com",
    "password": "secret123"
  }'
```

Login:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "fahri",
    "password": "secret123"
  }'
```

Buat kategori:

```bash
curl -X POST http://localhost:8080/api/v1/categories/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Makan",
    "description": "Pengeluaran makan harian"
  }'
```

Buat salary:

```bash
curl -X POST http://localhost:8080/api/v1/salary/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "5000000",
    "source": "Gaji utama",
    "description": "Salary bulanan"
  }'
```

Catatan:

- Field numerik seperti `amount` dan `allocated_amount` saat ini dikirim dalam bentuk string.
- ID pada endpoint detail/update/delete menggunakan ID bisnis seperti `salary_id`, `category_id`, atau `budget_id`, bukan integer primary key database.

## Pengembangan

Beberapa langkah yang biasanya dipakai saat development:

```bash
go fmt ./...
go test ./...
```

Kalau kamu ingin membuka project ini untuk contributor lain, bagus juga untuk menambahkan file `LICENSE` dan dokumentasi API yang lebih lengkap seperti Swagger atau OpenAPI.
