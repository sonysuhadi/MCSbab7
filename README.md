# MCS BAB 7 – Servo API

Aplikasi REST API sederhana untuk mengelola status servo menggunakan **Golang**, **Gin Framework**, dan **PostgreSQL**.

Project ini dibuat sebagai implementasi praktikum BAB 7 dengan konsep RESTful API dan pemisahan kode menggunakan arsitektur sederhana berbasis Controller, Repository, dan Entity.

## Teknologi yang Digunakan

- **Golang** – Bahasa pemrograman utama.
- **Gin Gonic** – Framework untuk membangun REST API.
- **PostgreSQL** – Database relasional.
- **lib/pq** – Driver PostgreSQL untuk Go.
- **sql-migrate** – Mengelola migrasi database.

## Struktur Project

```text
MCSbab7/
├── controllers/
│   └── servo_controller.go
├── database/
│   ├── sql_migrations/
│   │   └── 1_initiate.sql
│   └── database.go
├── entities/
│   └── servo_entity.go
├── repositories/
│   └── servo_repo.go
├── routers/
│   └── servo_router.go
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Persyaratan

Pastikan perangkat sudah memiliki:

- Go
- PostgreSQL
- Git
- Postman (untuk pengujian API)

## Instalasi dan Konfigurasi

### 1. Clone Repository

```bash
git clone https://github.com/Sonysuhadi/MCSbab7.git
```

Masuk ke direktori project:

```bash
cd MCSbab7
```

### 2. Instal Dependensi

```bash
go mod tidy
```

### 3. Buat Database

Buka PostgreSQL atau pgAdmin, kemudian buat database dengan nama:

```sql
CREATE DATABASE bab7;
```

### 4. Konfigurasi Database

Sesuaikan konfigurasi koneksi PostgreSQL di `main.go` dengan pengaturan PostgreSQL lokal:

- Host: localhost
- Port: 5432
- User: postgres
- Database: bab7

Pastikan password database sesuai dengan konfigurasi PostgreSQL milik masing-masing.

### 5. Jalankan Aplikasi

```bash
go run main.go
```

Jika berhasil, server berjalan pada:

```text
http://localhost:50001
```

Migrasi database akan dijalankan secara otomatis ketika aplikasi dimulai.

## API Endpoints

Base URL:

```text
http://localhost:50001
```

| Method | Endpoint | Fungsi |
|---|---|---|
| POST | `/servo/init-proj` | Inisialisasi data servo |
| GET | `/servo/status` | Mengambil status servo |
| PUT | `/servo/update/:srv_status` | Memperbarui status servo |

### 1. Inisialisasi Servo

```http
POST /servo/init-proj
```

Digunakan untuk membuat data awal servo dengan ID 1 dan status 0.

### 2. Mendapatkan Status Servo

```http
GET /servo/status
```

Digunakan untuk mengambil data status servo dari database.

### 3. Memperbarui Status Servo

```http
PUT /servo/update/90
```

Digunakan untuk memperbarui status servo menjadi 90.

Angka pada URL merupakan nilai status yang akan disimpan.

## Database

Tabel yang digunakan:

**Table: `servo`**

| Kolom | Tipe Data | Keterangan |
|---|---|---|
| id | INTEGER | Primary Key |
| status_servo | INTEGER | Nilai status servo |

## Pengujian API

Pengujian endpoint dapat dilakukan menggunakan Postman dengan metode HTTP:

- POST untuk inisialisasi.
- GET untuk mengambil status.
- PUT untuk memperbarui status.

## Developer

**Sony Kurniawan Suhadi**

GitHub: [@Sonysuhadi](https://github.com/Sonysuhadi)

Repository: [MCSbab7](https://github.com/Sonysuhadi/MCSbab7)

---

**MCS BAB 7 – Servo API**

Built with Go, Gin, and PostgreSQL.
