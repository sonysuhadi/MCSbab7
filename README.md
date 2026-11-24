# MCS BAB 7 – Membangun REST API Servo dengan Golang dan PostgreSQL

Tutorial lengkap membangun REST API untuk mengelola status servo menggunakan Golang, Gin Framework, dan PostgreSQL.

Tutorial ini dimulai dari instalasi tools, pembuatan project, penulisan kode setiap file, konfigurasi database, menjalankan server, hingga pengujian menggunakan Postman.

---

# DAFTAR ISI

1. Persiapan Tools
2. Instalasi Golang
3. Membuat Project
4. Membuat Struktur Folder
5. Inisialisasi Go Module
6. Instalasi Dependencies
7. Membuat Entity
8. Membuat Database dan Migrasi
9. Membuat Repository
10. Membuat Controller
11. Membuat Router
12. Membuat Main Program
13. Konfigurasi PostgreSQL
14. Menjalankan Aplikasi
15. Pengujian Menggunakan Postman
16. Memeriksa Database
17. Troubleshooting

---

# 1. Persiapan Tools

Sebelum memulai, siapkan beberapa aplikasi berikut:

| Tools | Kegunaan |
|---|---|
| Golang | Bahasa pemrograman |
| Visual Studio Code | Editor kode |
| PostgreSQL | Database |
| pgAdmin 4 | Mengelola database |
| Postman | Menguji REST API |
| Git | Version control |

Pastikan seluruh aplikasi sudah terinstal.

# 2. Instalasi Golang

Unduh Golang melalui:

https://go.dev/dl/

Setelah instalasi selesai, buka CMD atau PowerShell.

Periksa versi Golang:

```bash
go version
```

Jika berhasil, versi Golang akan ditampilkan.

# 3. Membuat Project

Buat folder project bernama `bab7`.

Contoh lokasi:

```text
C:\Users\sony\Documents\app flutter\bab7
```

Buka folder tersebut menggunakan Visual Studio Code.

Caranya:

1. Buka Visual Studio Code.
2. Pilih File → Open Folder.
3. Pilih folder `bab7`.
4. Buka terminal dengan Ctrl + `.

# 4. Membuat Struktur Folder

Buat folder dan file dengan struktur berikut:

```text
bab7/
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
└── main.go
```

Fungsi masing-masing folder:

- `controllers`: Mengatur request dan response API.
- `database`: Mengatur koneksi dan migrasi database.
- `entities`: Menentukan struktur data servo.
- `repositories`: Menjalankan query SQL.
- `routers`: Mengatur endpoint API.
- `main.go`: Program utama untuk menjalankan server.

# 5. Inisialisasi Go Module

Buka terminal pada folder project.

Jalankan:

```bash
go mod init bab7
```

Perintah ini membuat file `go.mod` sebagai identitas module Go.

# 6. Instalasi Dependencies

Jalankan:

```bash
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/rubenv/sql-migrate
```

Keterangan:

- Gin digunakan untuk membuat REST API.
- lib/pq digunakan sebagai driver PostgreSQL.
- sql-migrate digunakan untuk migrasi database.

# 7. Membuat Entity

Buka folder `entities`.

Buat file:

`servo_entity.go`

Isi kode:

```go
package entities

type Servo struct {
    Id          int `json:"id"`
    StatusServo int `json:"srv_status"`
}
```

Penjelasan:

- `Id` digunakan sebagai identitas data servo.
- `StatusServo` digunakan untuk menyimpan nilai status servo.
- JSON tag menentukan nama properti pada respons JSON.

# 8. Membuat Database dan Migrasi

## A. Membuat file database.go

Masuk ke folder `database`.

Buat file:

`database.go`

Isi kode:

```go
package database

import (
    "database/sql"
    "embed"
    "fmt"

    migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigrations embed.FS

var DbConnection *sql.DB

func DBMigrations(dbParam *sql.DB) error {
    migrations := &migrate.EmbedFileSystemMigrationSource{
        FileSystem: dbMigrations,
        Root:       "sql_migrations",
    }

    n, err := migrate.Exec(
        dbParam,
        "postgres",
        migrations,
        migrate.Up,
    )

    if err != nil {
        return fmt.Errorf("gagal menjalankan migrasi: %w", err)
    }

    DbConnection = dbParam

    fmt.Printf("Migrasi berhasil dijalankan: %d migrasi\n", n)

    return nil
}
```

Penjelasan:

- `embed.FS` digunakan untuk membaca file migrasi.
- `go:embed` menyertakan file SQL ke dalam aplikasi.
- `migrate.Exec` menjalankan migrasi database.
- `DbConnection` menyimpan koneksi database yang digunakan aplikasi.

## B. Membuat file migrasi SQL

Masuk ke:

`database/sql_migrations/`

Buat file:

`1_initiate.sql`

Isi:

```sql
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE servo (
    id INTEGER PRIMARY KEY,
    status_servo INTEGER
);

-- +migrate StatementEnd
```

Penjelasan:

- `CREATE TABLE servo` membuat tabel servo.
- `id` menjadi primary key.
- `status_servo` menyimpan nilai status servo.

Tabel yang dihasilkan:

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | INTEGER | Primary Key |
| status_servo | INTEGER | Status servo |

# 9. Membuat Repository

Masuk ke folder `repositories`.

Buat file:

`servo_repo.go`

Isi kode:

```go
package repositories

import (
    "bab7/entities"
    "database/sql"
)

func InitProj(db *sql.DB) (err error) {
    query := "INSERT INTO servo (id, status_servo) VALUES (1, 0) ON CONFLICT (id) DO NOTHING"
    _, err = db.Exec(query)
    return err
}

func GetStatus(db *sql.DB) (result []entities.Servo, err error) {
    query := "SELECT id, status_servo FROM servo"
    rows, err := db.Query(query)

    if err != nil {
        return
    }
    defer rows.Close()

    for rows.Next() {
        var data entities.Servo

        err = rows.Scan(&data.Id, &data.StatusServo)
        if err != nil {
            return
        }

        result = append(result, data)
    }

    err = rows.Err()
    return
}

func UpdateStatus(db *sql.DB, status entities.Servo) (err error) {
    query := "UPDATE servo SET status_servo = $1 WHERE id = 1"
    _, err = db.Exec(query, status.StatusServo)
    return
}
```

Penjelasan fungsi:

**InitProj**

Membuat data awal servo dengan ID 1 dan status 0. Jika data sudah ada, data tidak dibuat ulang.

**GetStatus**

Mengambil data servo dari PostgreSQL.

**UpdateStatus**

Memperbarui nilai status servo berdasarkan ID 1.

# 10. Membuat Controller

Masuk ke folder `controllers`.

Buat file:

`servo_controller.go`

Isi kode:

```go
package controllers

import (
    "bab7/database"
    "bab7/entities"
    "bab7/repositories"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

func InitProj(c *gin.Context) {
    err := repositories.InitProj(database.DbConnection)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{})
}

func GetStatus(c *gin.Context) {
    var result gin.H

    status, err := repositories.GetStatus(database.DbConnection)

    if err != nil {
        result = gin.H{
            "result": err.Error(),
        }
    } else {
        result = gin.H{
            "result": status,
        }
    }

    c.JSON(http.StatusOK, result)
}

func UpdateStatus(c *gin.Context) {
    var status entities.Servo

    srv_status, _ := strconv.Atoi(c.Param("srv_status"))
    status.StatusServo = srv_status

    err := repositories.UpdateStatus(database.DbConnection, status)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "StatusServo": status.StatusServo,
    })
}
```

Penjelasan:

- `InitProj`: Menginisialisasi data servo.
- `GetStatus`: Mengambil status servo.
- `UpdateStatus`: Mengubah status servo berdasarkan parameter URL.

# 11. Membuat Router

Masuk ke folder `routers`.

Buat file:

`servo_router.go`

Isi kode:

```go
package routers

import (
    "bab7/controllers"

    "github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
    routers := gin.Default()

    routers.POST(
        "/servo/init-proj",
        controllers.InitProj,
    )

    routers.GET(
        "/servo/status",
        controllers.GetStatus,
    )

    routers.PUT(
        "/servo/update/:srv_status",
        controllers.UpdateStatus,
    )

    return routers
}
```

Penjelasan:

Router digunakan untuk menghubungkan endpoint HTTP dengan fungsi controller.

# 12. Membuat Main Program

Buka file:

`main.go`

Isi dengan kode berikut:

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "bab7/database"
    "bab7/routers"

    _ "github.com/lib/pq"
)

func main() {
    const port = ":50001"

    host := "localhost"
    dbPort := 5432
    user := "postgres"
    dbName := "bab7"
    password := os.Getenv("DB_PASSWORD")

    if password == "" {
        log.Fatal("Environment variable DB_PASSWORD belum diatur")
    }

    psqlInfo := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host,
        dbPort,
        user,
        password,
        dbName,
    )

    DB, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Gagal membuka database: ", err)
    }
    defer DB.Close()

    err = DB.Ping()
    if err != nil {
        log.Fatal("Gagal menghubungkan PostgreSQL: ", err)
    }

    fmt.Println("DB berhasil tersambung")

    err = database.DBMigrations(DB)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Server berjalan di port ", port)

    err = routers.StartServer().Run(port)
    if err != nil {
        log.Fatal("Gagal menjalankan server: ", err)
    }
}
```

**Catatan:** Versi ini menggunakan environment variable agar password tidak ditulis langsung dalam source code. Pastikan mengganti `main.go` lokal dan memperbarui file tersebut di GitHub.

# 13. Konfigurasi PostgreSQL

Buka pgAdmin 4.

Buat database baru dengan nama:

```text
bab7
```

Atau melalui Query Tool:

```sql
CREATE DATABASE bab7;
```

Pastikan PostgreSQL berjalan pada port 5432.

## Mengatur password database di Windows

Jika menggunakan PowerShell:

```powershell
$env:DB_PASSWORD="PASSWORD_POSTGRESQL_KAMU"
```

Jika menggunakan CMD:

```cmd
set DB_PASSWORD=PASSWORD_POSTGRESQL_KAMU
```

Gunakan password PostgreSQL lokal milik masing-masing. Jangan memasukkan password ke README atau GitHub.

# 14. Menjalankan Aplikasi

Pastikan terminal berada di folder project.

Jalankan:

```bash
go mod tidy
```

Kemudian:

```bash
go run main.go
```

Jika berhasil, output kurang lebih:

```text
DB berhasil tersambung
Migrasi berhasil dijalankan: 1 migrasi
Server berjalan di port :50001

[GIN-debug] POST /servo/init-proj
[GIN-debug] GET /servo/status
[GIN-debug] PUT /servo/update/:srv_status

[GIN-debug] Listening and serving HTTP on :50001
```

Server sekarang aktif pada:

```text
http://localhost:50001
```

Biarkan terminal tetap terbuka selama pengujian.

# 15. Pengujian API Menggunakan Postman

Buka Postman dan buat request baru.

Base URL:

```text
http://localhost:50001
```

## A. Inisialisasi Data Servo

Method:

```text
POST
```

URL:

```text
http://localhost:50001/servo/init-proj
```

Klik Send.

Fungsi: Membuat data awal servo dengan ID 1 dan status 0.

## B. Melihat Status Servo

Method:

```text
GET
```

URL:

```text
http://localhost:50001/servo/status
```

Klik Send.

Contoh response:

```json
{
    "result": [
        {
            "id": 1,
            "srv_status": 0
        }
    ]
}
```

## C. Memperbarui Status Servo

Method:

```text
PUT
```

URL:

```text
http://localhost:50001/servo/update/90
```

Klik Send.

Contoh response:

```json
{
    "StatusServo": 90
}
```

Angka 90 merupakan contoh nilai status servo yang akan disimpan.

# 16. Memeriksa Database

Buka pgAdmin 4.

Pilih database `bab7`, kemudian buka Query Tool.

Jalankan:

```sql
SELECT * FROM servo;
```

Contoh hasil:

| id | status_servo |
|---|---|
| 1 | 90 |

Jika nilai status berubah menjadi 90, berarti proses UPDATE berhasil.

# 17. Alur Kerja Aplikasi

```text
        POSTMAN
           |
           v
      GIN ROUTER
           |
           v
       CONTROLLER
           |
           v
       REPOSITORY
           |
           v
       POSTGRESQL
           |
           v
      RESPONSE JSON
```

Penjelasan:

1. Postman mengirimkan HTTP request.
2. Gin Router menerima dan mengarahkan request.
3. Controller memproses request.
4. Repository menjalankan query SQL.
5. PostgreSQL menyimpan atau mengambil data.
6. API mengirimkan response dalam format JSON.

# 18. Troubleshooting

| Error | Penyebab dan solusi |
|---|---|
| `DB berhasil tersambung` tidak muncul | Periksa koneksi PostgreSQL dan password |
| `column statusservo does not exist` | Pastikan query menggunakan `status_servo` |
| `relation servo does not exist` | Pastikan migrasi database berhasil |
| `connection refused` | Pastikan PostgreSQL sedang berjalan |
| `port already in use` | Periksa penggunaan port 50001 |
| `go.mod` error | Jalankan `go mod tidy` |
| `DB_PASSWORD belum diatur` | Atur environment variable sebelum menjalankan aplikasi |

---

# Repository GitHub

Source code project:

https://github.com/Sonysuhadi/MCSbab7

Developer: **Sony Kurniawan Suhadi**

Project: **MCS BAB 7 – Servo API**

Built with Golang, Gin Framework, and PostgreSQL.
