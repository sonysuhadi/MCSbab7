# MCS BAB 7 – Servo API

REST API sederhana untuk mengelola status servo menggunakan Golang, Gin Framework, dan PostgreSQL.

Project ini dibuat untuk praktikum BAB 7 dengan menerapkan konsep RESTful API serta pemisahan kode menggunakan Controller, Repository, dan Entity.

---

## 1. Teknologi yang Digunakan

| Teknologi | Fungsi |
|---|---|
| Golang | Bahasa pemrograman |
| Gin Gonic | Framework REST API |
| PostgreSQL | Database |
| pgAdmin 4 | Manajemen database |
| lib/pq | Driver PostgreSQL |
| sql-migrate | Migrasi database |
| Postman | Pengujian API |

## 2. Struktur Project

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

## 3. Persiapan

Sebelum menjalankan project, pastikan sudah menginstal:

1. Go (Golang).
2. PostgreSQL.
3. pgAdmin 4.
4. Git.
5. Postman.
6. Visual Studio Code (opsional).

Periksa instalasi Go melalui terminal:

```bash
go version
```

Periksa Git:

```bash
git --version
```

## 4. Instalasi Project

### Langkah 1 – Clone Repository

Buka CMD atau terminal, kemudian jalankan:

```bash
git clone https://github.com/Sonysuhadi/MCSbab7.git
```

Masuk ke folder project:

```bash
cd MCSbab7
```

### Langkah 2 – Instal Dependensi

Jalankan perintah:

```bash
go mod tidy
```

Perintah tersebut digunakan untuk mengunduh dependensi yang dibutuhkan oleh project.

## 5. Konfigurasi Database PostgreSQL

### Langkah 1 – Membuka pgAdmin

1. Jalankan PostgreSQL.
2. Buka pgAdmin 4.
3. Hubungkan ke server PostgreSQL.
4. Klik kanan pada bagian Databases.
5. Pilih Create → Database.

### Langkah 2 – Membuat Database

Buat database dengan nama:

```text
bab7
```

Atau jalankan SQL berikut:

```sql
CREATE DATABASE bab7;
```

### Langkah 3 – Konfigurasi Koneksi

Buka file `main.go`.

Sesuaikan konfigurasi database dengan PostgreSQL yang digunakan:

| Parameter | Nilai |
|---|---|
| Host | localhost |
| Port | 5432 |
| User | postgres |
| Password | Password PostgreSQL masing-masing |
| Database | bab7 |

Pastikan konfigurasi koneksi sesuai dengan PostgreSQL lokal.

**Catatan keamanan:** Jangan mengunggah password database atau kredensial pribadi ke repository publik.

## 6. Menjalankan Aplikasi

Setelah konfigurasi database selesai, buka terminal pada direktori project.

Jalankan:

```bash
go run main.go
```

Jika berhasil, terminal akan menampilkan informasi seperti:

```text
DB berhasil tersambung
Migrasi berhasil dijalankan: 1 migrasi
Server berjalan di port :50001

[GIN-debug] POST   /servo/init-proj
[GIN-debug] GET    /servo/status
[GIN-debug] PUT    /servo/update/:srv_status

[GIN-debug] Listening and serving HTTP on :50001
```

Artinya:
- PostgreSQL berhasil terhubung.
- Migrasi database berhasil dijalankan.
- Endpoint API berhasil didaftarkan.
- Server berjalan pada port 50001.

Biarkan terminal tetap berjalan selama melakukan pengujian API.

## 7. Pengujian REST API Menggunakan Postman

Base URL:

```text
http://localhost:50001
```

### A. Inisialisasi Data Servo

Metode: POST

Endpoint:

```text
http://localhost:50001/servo/init-proj
```

Langkah pengujian:

1. Buka Postman.
2. Buat request baru.
3. Pilih metode POST.
4. Masukkan URL endpoint.
5. Klik Send.

Fungsi: Membuat data awal servo dengan ID 1 dan status 0.

Query database yang digunakan:

```sql
INSERT INTO servo (id, status_servo)
VALUES (1, 0)
ON CONFLICT (id) DO NOTHING;
```

### B. Melihat Status Servo

Metode: GET

Endpoint:

```text
http://localhost:50001/servo/status
```

Langkah pengujian:

1. Buat request baru.
2. Pilih metode GET.
3. Masukkan URL endpoint.
4. Klik Send.

Contoh respons:

```json
{
    "result": [
        {
            "Id": 1,
            "StatusServo": 0
        }
    ]
}
```

Fungsi: Mengambil data servo yang tersimpan di PostgreSQL.

### C. Memperbarui Status Servo

Metode: PUT

Endpoint:

```text
http://localhost:50001/servo/update/90
```

Langkah pengujian:

1. Buat request baru.
2. Pilih metode PUT.
3. Masukkan URL endpoint.
4. Klik Send.

Contoh respons:

```json
{
    "StatusServo": 90
}
```

Fungsi: Mengubah nilai status servo menjadi 90.

Angka 90 pada URL dapat diganti dengan nilai lain sesuai kebutuhan pengujian.

## 8. Memeriksa Data di PostgreSQL

Buka pgAdmin 4, kemudian pilih database `bab7`.

Buka Query Tool dan jalankan:

```sql
SELECT * FROM servo;
```

Contoh hasil:

| id | status_servo |
|---|---|
| 1 | 90 |

Nilai pada database akan berubah sesuai request PUT yang dikirim melalui Postman.

## 9. Alur Kerja Aplikasi

```text
Postman
   |
   v
Gin Router
   |
   v
Controller
   |
   v
Repository
   |
   v
PostgreSQL
   |
   v
Response JSON
```

Penjelasan:

1. Postman mengirimkan HTTP request.
2. Gin Router mengarahkan request ke controller.
3. Controller memproses parameter dan memanggil repository.
4. Repository menjalankan query SQL ke PostgreSQL.
5. Database mengembalikan hasil.
6. API mengirimkan respons dalam format JSON.

## 10. Troubleshooting

| Masalah | Solusi |
|---|---|
| Database connection failed | Periksa PostgreSQL, port, username, dan password |
| Relation servo does not exist | Pastikan migrasi database berhasil |
| Column statusservo does not exist | Gunakan nama kolom status_servo |
| Connection refused | Pastikan PostgreSQL sedang berjalan |
| Port already in use | Periksa apakah port 50001 digunakan aplikasi lain |
| Go module error | Jalankan `go mod tidy` |

## 11. GitHub

Repository: https://github.com/Sonysuhadi/MCSbab7

Developer: **Sony Kurniawan Suhadi**

---

**MCS BAB 7 – Servo API**

Built with Golang, Gin Framework, and PostgreSQL.
