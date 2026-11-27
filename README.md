# MCS BAB 7 – Kontrol Motor Servo Menggunakan ESP32, Go, dan PostgreSQL

## 1. Persiapan

Project ini menggunakan bahasa pemrograman Go sebagai REST API dan PostgreSQL sebagai database untuk menyimpan nilai sudut motor servo.

Pengujian dilakukan menggunakan Postman untuk mengirim perintah, membaca status, dan memperbarui nilai sudut servo.

### Software yang dibutuhkan

- Go
- PostgreSQL
- pgAdmin 4
- Visual Studio Code
- Postman

### Perangkat yang digunakan pada keseluruhan project

- ESP32
- Motor Servo SG90
- Kabel jumper
- Sumber daya 5V untuk servo

**Catatan:** Panduan ini hanya membahas backend Go, PostgreSQL, dan pengujian API menggunakan Postman. Pemrograman Arduino dan ESP32 tidak termasuk dalam pembahasan.

---

## 2. Membuat Database PostgreSQL

Pada tahap ini, kita akan membuat database menggunakan pgAdmin 4.

### Langkah 1 – Membuka pgAdmin 4

1. Buka aplikasi pgAdmin 4.
2. Pada bagian kiri, buka menu **Servers**.
3. Pilih server PostgreSQL yang digunakan.
4. Jika diminta, masukkan password PostgreSQL.

### Langkah 2 – Membuat Database

1. Klik kanan pada menu **Databases**.
2. Pilih **Create**.
3. Klik **Database**.
4. Masukkan nama database:

```text
bab7
```

5. Pada bagian Owner, pilih `postgres`.
6. Klik **Save**.

### Langkah 3 – Memastikan Database Berhasil Dibuat

Pada panel sebelah kiri pgAdmin, buka:

```text
Servers
└── PostgreSQL
    └── Databases
        └── bab7
```

Jika database `bab7` muncul, berarti proses pembuatan database berhasil.

### Langkah 4 – Membuka Query Tool

1. Klik database `bab7`.
2. Klik kanan pada database tersebut.
3. Pilih **Query Tool**.

Query Tool digunakan untuk menjalankan perintah SQL dan memeriksa isi database.

**Catatan:** Tabel `servo` tidak perlu dibuat secara manual karena akan dibuat otomatis oleh program Go melalui proses database migration.

---

## 3. Membuat Project Go

### Langkah 1 – Membuat Folder Project

Buka terminal atau PowerShell, kemudian jalankan:

```powershell
mkdir bab7
cd bab7
```

### Langkah 2 – Inisialisasi Project

Jalankan perintah:

```bash
go mod init bab7
```

Perintah tersebut digunakan untuk membuat file `go.mod` sebagai identitas module project Go.

### Langkah 3 – Menginstal Library

Jalankan perintah berikut:

```bash
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/rubenv/sql-migrate
```

Fungsi masing-masing library:

- `gin-gonic/gin`: Membuat REST API dan mengatur routing.
- `lib/pq`: Driver untuk menghubungkan Go dengan PostgreSQL.
- `sql-migrate`: Mengelola pembuatan dan perubahan struktur database melalui migration.

### Langkah 4 – Membuat Struktur Folder

Buat struktur folder dan file seperti berikut:

```text
bab7/
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
├── go.mod
├── go.sum
└── main.go
```

Setiap folder memiliki fungsi masing-masing untuk memisahkan proses pengelolaan data, koneksi database, pengaturan API, dan struktur data.

---

## 4. Memasukkan Kode Program Go

Masukkan kode berikut sesuai dengan nama file masing-masing.

### A. `entities/servo_entity.go`

```go
package entities

type Servo struct {
    Id          int `json:"id"`
    StatusServo int `json:"srv_status"`
}
```

**Fungsi:**

File ini digunakan untuk menentukan struktur data servo yang terdiri dari:

- `Id`: Identitas data servo.
- `StatusServo`: Nilai sudut motor servo.

Tag JSON digunakan agar nama field yang dikirim melalui API sesuai dengan format JSON.

### B. `database/database.go`

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

**Fungsi:**

File ini digunakan untuk menghubungkan proses migration dengan database PostgreSQL.

Penjelasan:

- `embed`: Memasukkan file SQL migration ke dalam program Go.
- `dbMigrations`: Menyimpan file migration.
- `DbConnection`: Menyimpan koneksi database yang digunakan oleh aplikasi.
- `migrate.Exec()`: Menjalankan proses migration secara otomatis.

Dengan demikian, tabel database dapat dibuat tanpa harus menjalankan perintah SQL secara manual setiap kali project dijalankan.

### C. `database/sql_migrations/1_initiate.sql`

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

**Fungsi:**

File ini digunakan untuk membuat tabel `servo` di PostgreSQL.

Penjelasan struktur tabel:

| Kolom | Tipe Data | Fungsi |
|---|---|---|
| id | INTEGER | Identitas data servo |
| status_servo | INTEGER | Menyimpan nilai sudut servo |

Keterangan tambahan:

- `PRIMARY KEY`: Menentukan kolom ID sebagai identitas utama.
- `DEFAULT 0`: Memberikan nilai awal sudut servo sebesar 0.
- `CHECK`: Membatasi nilai sudut antara 0 sampai 180 derajat.
- `Up`: Menjalankan pembuatan tabel.
- `Down`: Menghapus tabel ketika migration dibatalkan.

### D. `repositories/servo_repo.go`

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

**Fungsi:**

Repository merupakan bagian yang bertanggung jawab untuk berkomunikasi langsung dengan database.

Terdapat tiga fungsi utama:

1. `InitProj()`: Membuat data awal servo dengan ID 1 dan sudut 0. Jika data sudah tersedia, data tidak dibuat ulang.
2. `GetStatus()`: Mengambil data ID dan nilai sudut servo dari database.
3. `UpdateStatus()`: Memperbarui nilai sudut servo pada database.

### E. `controllers/servo_controller.go`

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

**Fungsi:**

Controller bertugas menerima permintaan dari pengguna, memanggil repository, dan mengirimkan hasil dalam format JSON.

Terdapat tiga fungsi:

- `InitProj()`: Memanggil proses pembuatan data awal.
- `GetStatus()`: Mengambil status servo dari database.
- `UpdateStatus()`: Menerima dan memperbarui nilai sudut servo.

Pada bagian `UpdateStatus()`, terdapat validasi agar nilai yang dikirim hanya berupa angka antara 0 sampai 180.

Jika nilai tidak sesuai, API akan mengembalikan pesan kesalahan dengan HTTP 400.

### F. `routers/servo_router.go`

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

**Fungsi:**

File router digunakan untuk menentukan alamat API dan menghubungkannya dengan controller.

Terdapat tiga endpoint:

| Method | Endpoint | Fungsi |
|---|---|---|
| POST | `/servo/init-proj` | Membuat data awal |
| GET | `/servo/status` | Melihat status servo |
| PUT | `/servo/update/:srv_status` | Memperbarui sudut servo |

### G. `main.go`

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "net"
    "net/url"
    "os"

    "bab7/database"
    "bab7/routers"

    _ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {

    value := os.Getenv(key)

    if value == "" {
        return fallback
    }

    return value
}

func main() {

    dbHost := getEnv("DB_HOST", "localhost")
    dbPort := getEnv("DB_PORT", "5432")
    dbUser := getEnv("DB_USER", "postgres")
    dbName := getEnv("DB_NAME", "bab7")
    dbPassword := os.Getenv("DB_PASSWORD")

    if dbPassword == "" {
        log.Fatal("Environment variable DB_PASSWORD belum diatur")
    }

    port := ":50001"

    dbURL := url.URL{
        Scheme: "postgres",
        User: url.UserPassword(
            dbUser,
            dbPassword,
        ),
        Host: net.JoinHostPort(dbHost, dbPort),
        Path: dbName,
    }

    query := dbURL.Query()
    query.Set("sslmode", "disable")
    dbURL.RawQuery = query.Encode()

    DB, err := sql.Open("postgres", dbURL.String())

    if err != nil {
        log.Fatal("Gagal membuka koneksi database:", err)
    }

    defer DB.Close()

    err = DB.Ping()

    if err != nil {
        log.Fatal("Database gagal tersambung:", err)
    }

    fmt.Println("DB berhasil tersambung")

    err = database.DBMigrations(DB)

    if err != nil {
        log.Fatal("Migrasi gagal:", err)
    }

    fmt.Println("Server berjalan di port", port)

    err = routers.StartServer().Run(port)

    if err != nil {
        log.Fatal("Server gagal dijalankan:", err)
    }
}
```

**Fungsi:**

File `main.go` merupakan program utama yang menjalankan seluruh proses aplikasi.

Tahapan yang dilakukan:

1. Membaca konfigurasi database dari environment variable.
2. Membuat koneksi ke PostgreSQL.
3. Memastikan koneksi database berhasil melalui `DB.Ping()`.
4. Menjalankan database migration.
5. Menjalankan REST API menggunakan Gin pada port `50001`.

**Penting:** Password PostgreSQL tidak ditulis langsung di dalam kode program. Password dibaca melalui environment variable `DB_PASSWORD` agar tidak mudah terekspos ketika project dibagikan atau diunggah ke GitHub.

---

## 5. Menjalankan dan Menguji Program Go

### Langkah 1 – Mengatur Password PostgreSQL

Sebelum menjalankan program, buka PowerShell pada folder project.

Masukkan perintah:

```powershell
$env:DB_PASSWORD="PASSWORD_POSTGRESQL_KAMU"
```

Ganti `PASSWORD_POSTGRESQL_KAMU` dengan password PostgreSQL yang digunakan.

Opsional, jika konfigurasi PostgreSQL tidak menggunakan nilai default:

```powershell
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_NAME="bab7"
```

Environment variable tersebut berlaku pada sesi PowerShell yang sedang digunakan.

### Langkah 2 – Memastikan Database Sudah Ada

Buka pgAdmin 4 dan pastikan database `bab7` sudah dibuat.

### Langkah 3 – Mengunduh Dependency

Jalankan perintah:

```bash
go mod tidy
```

Perintah ini digunakan untuk memastikan seluruh dependency yang dibutuhkan project sudah tersedia.

### Langkah 4 – Menjalankan Program

Jalankan:

```bash
go run .
```

Jika berhasil, akan muncul tampilan seperti berikut:

```text
DB berhasil tersambung
Migrasi berhasil dijalankan: 1 migrasi
Server berjalan di port :50001
```

Keterangan:

- `DB berhasil tersambung`: Go berhasil terhubung ke PostgreSQL.
- `Migrasi berhasil dijalankan`: Proses migration berhasil dijalankan.
- `Server berjalan di port :50001`: REST API sudah aktif dan siap menerima permintaan.

Biarkan terminal tetap berjalan selama pengujian API.

### Langkah 5 – Memeriksa Tabel di pgAdmin

Kembali ke pgAdmin 4, kemudian refresh database `bab7`.

Buka struktur berikut:

```text
bab7
└── Schemas
    └── public
        └── Tables
            └── servo
```

Jika tabel `servo` sudah muncul, berarti proses migration berhasil.

Untuk melihat struktur tabel, buka Query Tool dan jalankan:

```sql
SELECT * FROM servo;
```

Data awal belum tentu muncul karena data baru dibuat melalui endpoint inisialisasi pada tahap pengujian Postman.

---

## 6. Menguji REST API Menggunakan Postman

Pada tahap ini, Postman digunakan untuk mengirim permintaan ke REST API Go.

Pastikan server Go masih berjalan pada port `50001`.

### Pengujian 1 – Membuat Data Awal Servo

Buka Postman.

Atur konfigurasi:

- Method: `POST`
- URL: `http://localhost:50001/servo/init-proj`

Klik **Send**.

Contoh respons:

```json
{
    "message": "Data servo berhasil diinisialisasi"
}
```

Jika berhasil, data awal akan dibuat pada tabel `servo`.

Untuk memastikan data sudah masuk, buka Query Tool pada pgAdmin dan jalankan:

```sql
SELECT * FROM servo;
```

Contoh hasil:

| id | status_servo |
|---|---|
| 1 | 0 |

Artinya, data awal servo berhasil dibuat dengan sudut 0 derajat.

### Pengujian 2 – Melihat Status Servo

Pada Postman, atur:

- Method: `GET`
- URL: `http://localhost:50001/servo/status`

Klik **Send**.

Contoh respons:

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

Penjelasan:

- `id`: Identitas data servo.
- `srv_status`: Nilai sudut motor servo yang tersimpan dalam database.

Jika respons tersebut muncul, berarti Go berhasil membaca data dari PostgreSQL.

### Pengujian 3 – Mengubah Sudut Servo Menjadi 90 Derajat

Pada Postman, atur:

- Method: `PUT`
- URL: `http://localhost:50001/servo/update/90`

Klik **Send**.

Contoh respons:

```json
{
    "message": "Sudut servo berhasil diperbarui",
    "StatusServo": 90
}
```

Artinya, nilai sudut servo berhasil diperbarui menjadi 90 derajat.

### Pengujian 4 – Mengubah Sudut Servo Menjadi 180 Derajat

Atur URL:

```text
http://localhost:50001/servo/update/180
```

Gunakan method `PUT`, kemudian klik **Send**.

Contoh respons:

```json
{
    "message": "Sudut servo berhasil diperbarui",
    "StatusServo": 180
}
```

### Pengujian 5 – Mengubah Sudut Servo Menjadi 0 Derajat

Atur URL:

```text
http://localhost:50001/servo/update/0
```

Gunakan method `PUT`, kemudian klik **Send**.

Contoh respons:

```json
{
    "message": "Sudut servo berhasil diperbarui",
    "StatusServo": 0
}
```

### Pengujian 6 – Menguji Validasi Nilai Sudut

Untuk memastikan validasi berjalan, coba masukkan nilai di luar rentang yang diperbolehkan.

Contoh:

```text
PUT http://localhost:50001/servo/update/200
```

Hasil yang diharapkan:

- HTTP Status: `400 Bad Request`

Respons:

```json
{
    "error": "Sudut servo harus antara 0 sampai 180"
}
```

Coba juga memasukkan huruf:

```text
PUT http://localhost:50001/servo/update/abc
```

Hasil yang diharapkan:

```json
{
    "error": "Nilai sudut harus berupa angka"
}
```

Pengujian ini memastikan API hanya menerima nilai sudut berupa angka antara 0 sampai 180.

### Langkah 7 – Memastikan Perubahan Data di PostgreSQL

Setelah seluruh pengujian selesai, buka kembali Query Tool pada pgAdmin.

Jalankan:

```sql
SELECT * FROM servo;
```

Contoh hasil jika pengujian terakhir menggunakan sudut 0 derajat:

| id | status_servo |
|---|---|
| 1 | 0 |

Nilai `status_servo` akan mengikuti nilai terakhir yang berhasil dikirim melalui Postman.

---

## 7. Alur Kerja Sistem Backend

Berikut merupakan alur kerja sistem yang telah dibuat:

```text
             Postman
                |
                | Mengirim HTTP Request
                v
           Go REST API
                |
                v
            Controller
                |
                v
            Repository
                |
                | Query SQL
                v
           PostgreSQL
                |
                | Mengembalikan Data
                v
            Repository
                |
                v
            Controller
                |
                v
          JSON Response
                |
                v
             Postman
```

Penjelasan:

1. Postman mengirim permintaan HTTP ke server Go.
2. Router mengarahkan permintaan menuju controller yang sesuai.
3. Controller memanggil fungsi repository.
4. Repository menjalankan perintah SQL ke PostgreSQL.
5. PostgreSQL menyimpan atau mengembalikan data.
6. Controller mengirimkan hasil dalam format JSON.
7. Postman menampilkan respons dari server.

---

## 8. Kesimpulan

Pada project MCS BAB 7 ini, PostgreSQL digunakan sebagai database untuk menyimpan nilai sudut motor servo. Go digunakan sebagai REST API untuk mengelola komunikasi antara aplikasi dan database.

Dengan menggunakan framework Gin, sistem dapat menyediakan tiga endpoint utama, yaitu untuk inisialisasi data, membaca status servo, dan memperbarui nilai sudut servo.

Pengujian dilakukan menggunakan Postman untuk memastikan seluruh endpoint dapat berjalan dengan baik. Hasil perubahan data juga diperiksa melalui pgAdmin 4.

Dengan demikian, bagian backend telah berhasil dibuat dan diuji. Sistem sudah dapat menerima perintah, memproses data, serta menyimpan perubahan nilai sudut servo ke dalam PostgreSQL.

**Catatan:** Panduan ini berhenti pada tahap pengujian REST API menggunakan Postman. Integrasi ESP32 dan pemrograman Arduino merupakan tahap lanjutan yang tidak dibahas dalam dokumen ini.
