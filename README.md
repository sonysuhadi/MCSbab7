# MCS BAB 7 – Kontrol Motor Servo Menggunakan Go dan PostgreSQL

Project ini merupakan implementasi REST API menggunakan bahasa pemrograman Go (Golang) dan database PostgreSQL untuk mengelola data status motor servo.

Go digunakan untuk membuat REST API menggunakan framework Gin, sedangkan PostgreSQL digunakan untuk menyimpan data status servo. Pengujian dilakukan menggunakan Postman untuk mengirim perintah, membaca data, dan memperbarui status servo.

## 1. Teknologi yang Digunakan

- Go (Golang)
- Gin Framework
- PostgreSQL
- pgAdmin 4
- Visual Studio Code
- Postman
- sql-migrate
- lib/pq

## 2. Struktur Project

```text
MCSbab7/
│
├── controllers/
│   └── servo_controller.go
│
├── database/
│   ├── sql_migrations/
│   │   └── 1_initiate.sql
│   └── database.go
│
├── entities/
│   └── servo_entity.go
│
├── repositories/
│   └── servo_repo.go
│
├── routers/
│   └── servo_router.go
│
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## 3. Membuat Database PostgreSQL

1. Buka aplikasi pgAdmin 4.
2. Masuk ke menu Servers.
3. Pilih server PostgreSQL.
4. Klik kanan Databases.
5. Pilih Create → Database.
6. Masukkan nama database:

```text
bab7
```

7. Pada bagian Owner, pilih `postgres`.
8. Klik Save.

Tabel `servo` akan dibuat secara otomatis melalui proses database migration ketika program Go dijalankan.

## 4. Inisialisasi Project Go

Buka terminal pada folder project.

### A. Inisialisasi Module

```bash
go mod init bab7
```

### B. Instalasi Library

```bash
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/rubenv/sql-migrate
go get -u "github.com/gin-contrib/cors"
```

## 5. Kode Program

Masukkan kode berikut sesuai dengan struktur folder project.

### A. entities/servo_entity.go

```go
package entities

type Servo struct {
    Id          int `json:"id"`
    StatusServo int `json:"srv_status"`
}
```

Fungsi: Menentukan struktur data servo yang terdiri dari ID dan nilai sudut servo.

### B. database/database.go

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

Fungsi:
- Menghubungkan aplikasi dengan database PostgreSQL.
- Membaca file SQL migration.
- Menjalankan migration secara otomatis.
- Menyimpan koneksi database pada `DbConnection`.

### C. database/sql_migrations/1_initiate.sql

```sql
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE servo (
    id INTEGER PRIMARY KEY,
    status_servo INTEGER NOT NULL DEFAULT 0
        CHECK (status_servo >= 0 AND status_servo <= 180)
);

-- +migrate StatementEnd


-- +migrate Down
-- +migrate StatementBegin

DROP TABLE servo;

-- +migrate StatementEnd
```

Fungsi: Membuat tabel `servo` pada PostgreSQL.

| Kolom | Tipe Data | Fungsi |
|---|---|---|
| id | INTEGER | Identitas data servo |
| status_servo | INTEGER | Menyimpan nilai sudut servo |

### D. repositories/servo_repo.go

```go
package repositories

import (
    "bab7/entities"
    "database/sql"
    "errors"
)

func InitProj(db *sql.DB) error {

    query := `
        INSERT INTO servo (id, status_servo)
        VALUES (1, 0)
        ON CONFLICT (id) DO NOTHING
    `

    _, err := db.Exec(query)

    return err
}

func GetStatus(db *sql.DB) ([]entities.Servo, error) {

    query := `
        SELECT id, status_servo
        FROM servo
        ORDER BY id
    `

    rows, err := db.Query(query)

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var result []entities.Servo

    for rows.Next() {

        var data entities.Servo

        err = rows.Scan(
            &data.Id,
            &data.StatusServo,
        )

        if err != nil {
            return nil, err
        }

        result = append(result, data)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return result, nil
}

func UpdateStatus(db *sql.DB, status entities.Servo) error {

    query := `
        UPDATE servo
        SET status_servo = $1
        WHERE id = 1
    `

    result, err := db.Exec(query, status.StatusServo)

    if err != nil {
        return err
    }

    affected, err := result.RowsAffected()

    if err != nil {
        return err
    }

    if affected == 0 {
        return errors.New("data servo belum diinisialisasi")
    }

    return nil
}
```

Fungsi:
- `InitProj()`: Membuat data awal servo dengan ID 1 dan nilai 0.
- `GetStatus()`: Mengambil data servo dari PostgreSQL.
- `UpdateStatus()`: Memperbarui nilai sudut servo.

### E. controllers/servo_controller.go

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

    c.JSON(http.StatusOK, gin.H{
        "message": "Data servo berhasil diinisialisasi",
    })
}

func GetStatus(c *gin.Context) {

    status, err := repositories.GetStatus(database.DbConnection)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "result": status,
    })
}

func UpdateStatus(c *gin.Context) {

    srvStatus, err := strconv.Atoi(c.Param("srv_status"))

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Nilai sudut harus berupa angka",
        })
        return
    }

    if srvStatus < 0 || srvStatus > 180 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Sudut servo harus antara 0 sampai 180",
        })
        return
    }

    status := entities.Servo{
        StatusServo: srvStatus,
    }

    err = repositories.UpdateStatus(
        database.DbConnection,
        status,
    )

    if err != nil {

        if err.Error() == "data servo belum diinisialisasi" {
            c.JSON(http.StatusNotFound, gin.H{
                "error": err.Error(),
            })
            return
        }

        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":      "Sudut servo berhasil diperbarui",
        "StatusServo": status.StatusServo,
    })
}
```

Fungsi:
- `InitProj()`: Menginisialisasi data servo.
- `GetStatus()`: Mengambil status servo.
- `UpdateStatus()`: Memperbarui status servo berdasarkan parameter URL.

### F. routers/servo_router.go

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

Fungsi: Mengatur endpoint REST API dan menghubungkannya dengan controller.

Endpoint yang tersedia:

| Method | Endpoint | Fungsi |
|---|---|---|
| POST | `/servo/init-proj` | Membuat data awal |
| GET | `/servo/status` | Melihat status servo |
| PUT | `/servo/update/:srv_status` | Memperbarui sudut servo |

### G. main.go

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    "bab7/database"
    "bab7/routers"

    _ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "PASSWORD_POSTGRESQL_KAMU"
    dbName   = "bab7"
)

func main() {

    PORT := ":50001"

    psqlInfo := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host,
        port,
        user,
        password,
        dbName,
    )

    // Membuka koneksi database
    DB, err := sql.Open("postgres", psqlInfo)

    if err != nil {
        log.Fatal("Gagal membuka database: ", err)
    }

    defer DB.Close()

    // Memeriksa koneksi PostgreSQL
    err = DB.Ping()

    if err != nil {
        log.Fatal("Gagal menghubungkan PostgreSQL: ", err)
    }

    fmt.Println("DB berhasil tersambung")

    // Menjalankan migrasi database
    err = database.DBMigrations(DB)

    if err != nil {
        log.Fatal(err)
    }

    // Menjalankan server Gin
    fmt.Println("Server berjalan di port ", PORT)

    err = routers.StartServer().Run(PORT)

    if err != nil {
        log.Fatal("Gagal menjalankan server: ", err)
    }
}
```

Fungsi:
1. Mengatur konfigurasi koneksi PostgreSQL.
2. Membuka koneksi database.
3. Memeriksa koneksi menggunakan `DB.Ping()`.
4. Menjalankan database migration.
5. Menjalankan REST API menggunakan Gin pada port `50001`.

**Catatan:** Ganti `PASSWORD_POSTGRESQL_KAMU` dengan password PostgreSQL lokal. Jangan memasukkan password asli ke repository publik.

## 6. Menjalankan Program

Pastikan PostgreSQL sudah aktif dan database `bab7` sudah tersedia.

### A. Mengunduh Dependency

```bash
go mod tidy
```

### B. Menjalankan Server

```bash
go run main.go
```

Jika berhasil, terminal akan menampilkan:

```text
DB berhasil tersambung
Migrasi berhasil dijalankan: 1 migrasi
Server berjalan di port :50001
```

Server REST API berjalan pada:

```text
http://localhost:50001
```

Biarkan terminal tetap berjalan selama pengujian Postman.

## 7. Pengujian REST API Menggunakan Postman

### A. POST – Inisialisasi Data Servo

Method:

```http
POST
```

URL:

```text
http://localhost:50001/servo/init-proj
```

Contoh response:

```json
{
    "message": "Data servo berhasil diinisialisasi"
}
```

### B. GET – Melihat Status Servo

Method:

```http
GET
```

URL:

```text
http://localhost:50001/servo/status
```

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

### C. PUT – Mengubah Status Servo

Method:

```http
PUT
```

URL:

```text
http://localhost:50001/servo/update/1
```

Contoh response:

```json
{
    "message": "Sudut servo berhasil diperbarui",
    "StatusServo": 1
}
```

### D. GET – Memastikan Perubahan Data

Setelah melakukan PUT, jalankan kembali:

```text
http://localhost:50001/servo/status
```

Contoh response:

```json
{
    "result": [
        {
            "id": 1,
            "srv_status": 1
        }
    ]
}
```

Jika nilai `srv_status` berubah menjadi 1, berarti proses pembaruan data berhasil.

## 8. Memeriksa Data PostgreSQL

Buka pgAdmin 4, pilih database `bab7`, kemudian buka Query Tool.

Jalankan SQL:

```sql
SELECT * FROM servo;
```

Contoh hasil:

| id | status_servo |
|---|---|
| 1 | 1 |

Nilai tersebut menunjukkan bahwa data berhasil diperbarui melalui REST API dan tersimpan di PostgreSQL.

## 9. Hasil Pengujian

Berdasarkan pengujian menggunakan Postman, diperoleh hasil sebagai berikut:

| No. | Method | Endpoint | Hasil |
|---|---|---|---|
| 1 | POST | `/servo/init-proj` | Berhasil |
| 2 | GET | `/servo/status` | Berhasil |
| 3 | PUT | `/servo/update/1` | Berhasil |
| 4 | GET | `/servo/status` | Data berhasil diperbarui |

Hasil pengujian membuktikan bahwa:
- Go berhasil terhubung ke PostgreSQL.
- Database migration berhasil dijalankan.
- Server Gin berjalan pada port 50001.
- Data servo berhasil diinisialisasi.
- Data servo berhasil dibaca.
- Data servo berhasil diperbarui.
- Perubahan data tersimpan di PostgreSQL.

## 10. Kesimpulan

Project MCS BAB 7 berhasil mengimplementasikan REST API menggunakan Go dan PostgreSQL untuk mengelola data status motor servo.

Go digunakan untuk menangani request HTTP melalui framework Gin, sedangkan PostgreSQL digunakan untuk menyimpan data. Seluruh endpoint telah diuji menggunakan Postman untuk memastikan proses inisialisasi, pembacaan, dan pembaruan data berjalan dengan baik.

Dengan demikian, implementasi backend dan pengujian REST API telah selesai.

---

**Repository:** https://github.com/sonysuhadi/MCSbab7

**MCS BAB 7 – Go REST API & PostgreSQL**
