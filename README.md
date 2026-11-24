# MCS BAB 7 – Kontrol Motor Servo Menggunakan ESP32, Go, dan PostgreSQL

Project ini menggunakan ESP32 untuk mengontrol motor servo melalui REST API Go. Nilai sudut servo disimpan di PostgreSQL dan dibaca oleh ESP32 secara berkala.

## 1. Persiapan

Pastikan software dan perangkat berikut sudah tersedia:

**Software:**
- Go
- PostgreSQL
- pgAdmin 4
- Visual Studio Code
- Postman
- Arduino IDE

**Perangkat:**
- ESP32
- Motor Servo SG90
- Kabel jumper
- Sumber daya 5V untuk servo

## 2. Membuat Database PostgreSQL

Pada tahap ini kita membuat database menggunakan pgAdmin 4.

### Langkah 1 – Membuka pgAdmin 4

1. Buka aplikasi pgAdmin 4.
2. Pada bagian kiri, buka menu **Servers**.
3. Pilih server PostgreSQL yang digunakan.
4. Jika diminta, masukkan password PostgreSQL.

### Langkah 2 – Membuat Database

1. Klik kanan pada menu **Databases**.
2. Pilih **Create**.
3. Klik **Database**.
4. Isi nama database dengan:

```text
bab7
```

5. Pada bagian Owner, pilih `postgres`.
6. Klik **Save**.

Database `bab7` sekarang sudah dibuat.

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

Query Tool digunakan untuk menjalankan perintah SQL.

> Tabel `servo` tidak perlu dibuat secara manual karena akan dibuat otomatis oleh program Go melalui proses migrasi.

## 3. Membuat Project Go

### Langkah 1 – Membuat Folder Project

Buka terminal atau PowerShell, lalu jalankan:

```powershell
mkdir bab7
cd bab7
```

### Langkah 2 – Inisialisasi Project

```bash
go mod init bab7
```

### Langkah 3 – Install Library

```bash
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/rubenv/sql-migrate
```

### Langkah 4 – Membuat Struktur Folder

Buat folder dan file berikut:

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

## 4. Memasukkan Kode Program Go

Masukkan kode berikut sesuai nama file masing-masing.

### A. `entities/servo_entity.go`

```go
package entities

type Servo struct {
    Id          int `json:"id"`
    StatusServo int `json:"srv_status"`
}
```

Fungsi: Menentukan struktur data servo yang terdiri dari ID dan nilai sudut servo.

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

Fungsi: Menghubungkan proses migrasi dengan database PostgreSQL. File SQL yang ada di folder `sql_migrations` akan dijalankan secara otomatis.

### C. `database/sql_migrations/1_initiate.sql`

```sql
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE servo (
    id INTEGER PRIMARY KEY,
    status_servo INTEGER
);

-- +migrate StatementEnd
```

Fungsi: Membuat tabel `servo` untuk menyimpan ID dan nilai sudut motor servo.

### D. `repositories/servo_repo.go`

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

Fungsi:
- `InitProj()`: Membuat data awal servo dengan sudut 0.
- `GetStatus()`: Mengambil nilai sudut servo dari database.
- `UpdateStatus()`: Memperbarui nilai sudut servo di database.

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
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{})
}

func GetStatus(c *gin.Context) {
    var result gin.H

    status, err := repositories.GetStatus(database.DbConnection)

    if err != nil {
        result = gin.H{"result": err.Error()}
    } else {
        result = gin.H{"result": status}
    }

    c.JSON(http.StatusOK, result)
}

func UpdateStatus(c *gin.Context) {
    var status entities.Servo

    srv_status, _ := strconv.Atoi(c.Param("srv_status"))
    status.StatusServo = srv_status

    err := repositories.UpdateStatus(database.DbConnection, status)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"StatusServo": status.StatusServo})
}
```

Fungsi: Mengatur permintaan dari pengguna, memanggil repository, dan mengirimkan hasil dalam format JSON.

### F. `routers/servo_router.go`

```go
package routers

import (
    "bab7/controllers"

    "github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
    routers := gin.Default()

    routers.POST("/servo/init-proj", controllers.InitProj)
    routers.GET("/servo/status", controllers.GetStatus)
    routers.PUT("/servo/update/:srv_status", controllers.UpdateStatus)

    return routers
}
```

Fungsi: Menentukan tiga alamat API untuk membuat data awal, membaca status, dan memperbarui sudut servo.

### G. `main.go`

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

    DB, err := sql.Open("postgres", psqlInfo)

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

    fmt.Println("Server berjalan di port ", PORT)

    err = routers.StartServer().Run(PORT)

    if err != nil {
        log.Fatal("Server gagal dijalankan:", err)
    }
}
```

**Penting:** Ganti `PASSWORD_POSTGRESQL_KAMU` dengan password PostgreSQL milik sendiri. Jangan masukkan password asli ke GitHub publik. Kode ini adalah contoh konfigurasi untuk README; sesuaikan dengan konfigurasi `main.go` yang digunakan pada project.

## 5. Menjalankan dan Menguji API

### Langkah 1 – Memastikan Database Sudah Ada

Buka pgAdmin 4 dan pastikan database `bab7` sudah muncul.

### Langkah 2 – Menjalankan Program

Buka terminal di folder project, kemudian jalankan:

```bash
go mod tidy
go run main.go
```

Jika berhasil, akan muncul:

```text
DB berhasil tersambung
Migrasi berhasil dijalankan: 1 migrasi
Server berjalan di port :50001
```

Artinya:
- Go berhasil terhubung ke PostgreSQL.
- Migrasi berhasil membuat tabel `servo`.
- Server API berjalan pada port `50001`.

### Langkah 3 – Memeriksa Tabel di pgAdmin

Kembali ke pgAdmin 4, kemudian refresh database `bab7`.

Buka:

```text
bab7
└── Schemas
    └── public
        └── Tables
            └── servo
```

Jika tabel `servo` sudah muncul, berarti migrasi berhasil.

Untuk melihat isi tabel, buka Query Tool dan jalankan:

```sql
SELECT * FROM servo;
```

Data awal akan dibuat pada tahap pengujian API.

### Langkah 4 – Menguji API dengan Postman

Pastikan server Go masih berjalan.

**Membuat data awal**

- Method: `POST`
- URL: `http://localhost:50001/servo/init-proj`

**Melihat status servo**

- Method: `GET`
- URL: `http://localhost:50001/servo/status`

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

**Mengubah sudut servo**

- Method: `PUT`
- URL: `http://localhost:50001/servo/update/90`

Contoh respons:

```json
{
    "StatusServo": 90
}
```

Coba juga nilai `0`, `45`, dan `180`.

Setelah mengubah nilai, jalankan kembali:

```sql
SELECT * FROM servo;
```

Nilai `status_servo` seharusnya mengikuti nilai terakhir yang dikirim melalui Postman.

## 6. Menghubungkan ESP32 dengan Motor Servo

### Langkah 1 – Merangkai Komponen

Hubungkan motor servo SG90 ke ESP32:

| Kabel Servo | Sambungan |
|---|---|
| Merah (VCC) | Sumber daya 5V |
| Cokelat/Hitam (GND) | GND |
| Oranye/Kuning (Signal) | GPIO 13 |

Gunakan sumber daya 5V yang sesuai untuk servo. Satukan GND sumber daya servo dengan GND ESP32.

### Langkah 2 – Menginstal Library Arduino

Buka Arduino IDE, pilih **Library Manager**, lalu instal:

- ESP32Servo
- ArduinoJson

Library WiFi dan HTTPClient sudah tersedia melalui paket board ESP32.

### Langkah 3 – Mengetahui IP Address Komputer

ESP32 harus mengakses IP komputer tempat server Go berjalan.

Buka Command Prompt, lalu ketik:

```bash
ipconfig
```

Cari bagian `IPv4 Address` pada koneksi WiFi atau Ethernet yang sedang digunakan.

Contoh:

```text
IPv4 Address: 192.168.1.50
```

Pastikan komputer dan ESP32 menggunakan jaringan yang sama.

### Langkah 4 – Program ESP32

Masukkan kode berikut ke Arduino IDE.

Ganti nama WiFi, password WiFi, dan IP server sesuai jaringan yang digunakan.

```cpp
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include <ESP32Servo.h>

const char* ssid = "NAMA_WIFI_ANDA";
const char* password = "PASSWORD_WIFI_ANDA";

const char* serverUrl =
    "http://192.168.1.50:50001/servo/status";

Servo myServo;

const int servoPin = 13;
int lastAngle = -1;

void setup() {
    Serial.begin(115200);

    myServo.setPeriodHertz(50);
    myServo.attach(servoPin, 500, 2400);

    WiFi.begin(ssid, password);

    Serial.print("Menghubungkan WiFi");

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    Serial.println();
    Serial.println("WiFi berhasil terhubung");

    Serial.print("IP ESP32: ");
    Serial.println(WiFi.localIP());
}

void loop() {
    if (WiFi.status() == WL_CONNECTED) {

        HTTPClient http;

        http.begin(serverUrl);
        http.setTimeout(5000);

        int httpCode = http.GET();

        if (httpCode == HTTP_CODE_OK) {

            String payload = http.getString();

            Serial.println("Data dari server:");
            Serial.println(payload);

            DynamicJsonDocument doc(1024);

            DeserializationError error =
                deserializeJson(doc, payload);

            if (!error) {

                int targetAngle =
                    doc["result"][0]["srv_status"] | 0;

                targetAngle = constrain(targetAngle, 0, 180);

                if (targetAngle != lastAngle) {

                    myServo.write(targetAngle);

                    lastAngle = targetAngle;

                    Serial.print("Sudut servo: ");
                    Serial.println(targetAngle);
                }

            } else {
                Serial.println("Gagal membaca JSON");
            }

        } else {
            Serial.print("HTTP Error: ");
            Serial.println(httpCode);
        }

        http.end();

    } else {
        Serial.println("WiFi terputus");
    }

    delay(2000);
}
```

### Langkah 5 – Upload Program

1. Hubungkan ESP32 ke komputer.
2. Pilih board ESP32 di Arduino IDE.
3. Pilih port ESP32.
4. Klik Upload.
5. Setelah selesai, buka Serial Monitor.
6. Atur baud rate menjadi `115200`.

Jika berhasil, akan muncul pesan:

```text
Menghubungkan WiFi...
WiFi berhasil terhubung
IP ESP32: 192.168.1.100
```

## 7. Pengujian Akhir

Pastikan server Go masih berjalan dan ESP32 sudah terhubung ke WiFi.

### Pengujian 1 – Menggerakkan Servo ke 180 Derajat

Buka Postman, kemudian kirim:

```http
PUT http://localhost:50001/servo/update/180
```

ESP32 akan membaca data terbaru dari API dan menggerakkan servo menuju 180 derajat.

### Pengujian 2 – Menggerakkan Servo ke 0 Derajat

```http
PUT http://localhost:50001/servo/update/0
```

Servo akan bergerak menuju 0 derajat.

### Pengujian 3 – Menggerakkan Servo ke 90 Derajat

```http
PUT http://localhost:50001/servo/update/90
```

Servo akan bergerak menuju 90 derajat.

ESP32 membaca data dari server setiap 2 detik, sehingga perubahan sudut tidak harus dilakukan dengan mengunggah ulang program.

### Alur Kerja Sistem

```text
Postman
   |
   | Mengirim sudut servo
   v
Go REST API
   |
   | Menyimpan data
   v
PostgreSQL
   |
   | Data dibaca melalui API
   v
ESP32
   |
   | Mengatur sudut servo
   v
Motor Servo
```

**Kesimpulan**

Pada project ini, PostgreSQL digunakan untuk menyimpan nilai sudut motor servo. Go digunakan sebagai REST API untuk membaca dan memperbarui data. ESP32 mengambil data dari API melalui koneksi WiFi, kemudian menggerakkan motor servo sesuai nilai sudut yang tersimpan.

Dengan demikian, perubahan sudut motor servo dapat dilakukan melalui Postman tanpa perlu mengunggah ulang program ESP32.

**Catatan:** Gunakan nilai sudut 0–180 saat pengujian. Program backend contoh ini belum memvalidasi rentang nilai tersebut, sedangkan program ESP32 membatasi perintah gerakan dengan `constrain()`.
