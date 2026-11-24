# MCS BAB 7 – Monitoring dan Kontrol Motor Servo Menggunakan ESP32, Go, dan PostgreSQL

## 1. Persiapan

Software dan perangkat yang dibutuhkan:

- Go
- PostgreSQL
- Visual Studio Code
- Postman
- Arduino IDE
- ESP32
- Motor Servo SG90
- Kabel jumper

## 2. Struktur Folder Project

Buat project dengan nama `bab7`.

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

## 3. Membuat Program Go

### A. Inisialisasi Project

Buka terminal di folder project, lalu jalankan:

```bash
go mod init bab7
```

Install library yang dibutuhkan:

```bash
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/rubenv/sql-migrate
```

### B. Entity – `entities/servo_entity.go`

File ini digunakan untuk menentukan struktur data servo.

```go
package entities

type Servo struct {
    Id          int `json:"id"`
    StatusServo int `json:"srv_status"`
}
```

Penjelasan:
- `Id`: ID data servo.
- `StatusServo`: nilai sudut motor servo.
- `json`: menentukan nama field ketika data dikirim dalam format JSON.

### C. Database – `database/database.go`

File ini digunakan untuk menjalankan migrasi database.

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
- `embed.FS`: membaca file migrasi yang ada di dalam project.
- `migrate.Exec`: menjalankan migrasi database.
- `migrate.Up`: menjalankan migrasi untuk membuat tabel.
- `DbConnection`: menyimpan koneksi database agar dapat digunakan file lain.

### D. Migrasi Database – `database/sql_migrations/1_initiate.sql`

File ini digunakan untuk membuat tabel servo.

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
- `id`: identitas data servo.
- `status_servo`: menyimpan nilai sudut servo.
- `PRIMARY KEY`: menjadikan ID sebagai identitas unik.

### E. Repository – `repositories/servo_repo.go`

Repository berisi fungsi untuk berinteraksi langsung dengan PostgreSQL.

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

| Fungsi | Kegunaan |
|---|---|
| `InitProj()` | Membuat data awal servo dengan sudut 0. |
| `GetStatus()` | Mengambil data sudut servo dari database. |
| `UpdateStatus()` | Mengubah nilai sudut servo di database. |

### F. Controller – `controllers/servo_controller.go`

Controller menerima permintaan API dan menghubungkannya dengan repository.

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

Penjelasan fungsi:

| Fungsi | Kegunaan |
|---|---|
| `InitProj()` | Memanggil fungsi untuk membuat data awal. |
| `GetStatus()` | Mengirim data servo dari database dalam format JSON. |
| `UpdateStatus()` | Menerima sudut servo dari URL dan menyimpan nilainya. |

### G. Router – `routers/servo_router.go`

Router menentukan alamat API yang dapat diakses.

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

Terdapat tiga endpoint:

| Method | Endpoint | Fungsi |
|---|---|---|
| POST | `/servo/init-proj` | Membuat data awal |
| GET | `/servo/status` | Membaca status servo |
| PUT | `/servo/update/:srv_status` | Mengubah sudut servo |

### H. Main – `main.go`

File utama untuk menjalankan koneksi database, migrasi, dan server API.

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

**Catatan:** Ganti `PASSWORD_POSTGRESQL_KAMU` dengan password PostgreSQL milik sendiri. Jangan memasukkan password asli ke repository publik GitHub.

## 4. Menyiapkan PostgreSQL dan Menjalankan API

### A. Membuat Database

Buka pgAdmin atau PostgreSQL, lalu buat database dengan nama:

```sql
CREATE DATABASE bab7;
```

Pastikan konfigurasi database pada `main.go` sesuai dengan PostgreSQL yang digunakan.

### B. Menjalankan Program

Buka terminal di folder project:

```bash
go mod tidy
go run main.go
```

Jika berhasil, terminal akan menampilkan:

```text
DB berhasil tersambung
Migrasi berhasil dijalankan: 1 migrasi
Server berjalan di port :50001
```

Artinya koneksi database dan server API berhasil dijalankan.

### C. Menguji API dengan Postman

Pastikan server Go masih berjalan.

**1. Membuat data awal**

- Method: `POST`
- URL: `http://localhost:50001/servo/init-proj`

**2. Melihat status servo**

- Method: `GET`
- URL: `http://localhost:50001/servo/status`

Contoh hasil:

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

**3. Mengubah sudut servo**

- Method: `PUT`
- URL: `http://localhost:50001/servo/update/90`

Contoh hasil:

```json
{
    "StatusServo": 90
}
```

Coba juga nilai `0`, `45`, dan `180`.

## 5. Menghubungkan ESP32 dengan Motor Servo

### A. Rangkaian

Contoh koneksi servo SG90 ke ESP32:

| Kabel Servo | ESP32 / Sumber Daya |
|---|---|
| Merah (VCC) | Sumber 5V yang sesuai |
| Cokelat/Hitam (GND) | GND |
| Oranye/Kuning (Signal) | GPIO 13 |

Gunakan sumber daya 5V yang mampu menyediakan arus untuk servo. Satukan GND sumber daya servo dengan GND ESP32.

### B. Library Arduino IDE

Install library berikut melalui Library Manager:

- ESP32Servo
- ArduinoJson
- HTTPClient dan WiFi tersedia melalui paket board ESP32.

### C. Program ESP32 Lengkap

Ganti nama WiFi, password WiFi, dan alamat IP server sesuai jaringan yang digunakan.

```cpp
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include <ESP32Servo.h>

// Konfigurasi WiFi
const char* ssid = "NAMA_WIFI_ANDA";
const char* password = "PASSWORD_WIFI_ANDA";

// Alamat API Go
const char* serverUrl =
    "http://192.168.1.50:50001/servo/status";

// Konfigurasi servo
Servo myServo;

const int servoPin = 13;
int lastAngle = -1;

void setup() {
    Serial.begin(115200);

    // Mengatur pin servo
    myServo.setPeriodHertz(50);
    myServo.attach(servoPin, 500, 2400);

    // Menghubungkan ESP32 ke WiFi
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
    // Memastikan ESP32 terhubung ke WiFi
    if (WiFi.status() == WL_CONNECTED) {

        HTTPClient http;

        http.begin(serverUrl);
        http.setTimeout(5000);

        // Mengambil data dari API
        int httpCode = http.GET();

        if (httpCode == HTTP_CODE_OK) {

            String payload = http.getString();

            Serial.println("Data dari server:");
            Serial.println(payload);

            // Membaca JSON
            DynamicJsonDocument doc(1024);

            DeserializationError error =
                deserializeJson(doc, payload);

            if (!error) {

                // Mengambil nilai sudut servo
                int targetAngle =
                    doc["result"][0]["srv_status"] | 0;

                // Membatasi sudut 0 sampai 180
                targetAngle = constrain(targetAngle, 0, 180);

                // Menggerakkan servo jika nilai berubah
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

    // Memeriksa server setiap 2 detik
    delay(2000);
}
```

Penjelasan fungsi utama:

| Fungsi | Kegunaan |
|---|---|
| `WiFi.begin()` | Menghubungkan ESP32 ke WiFi. |
| `HTTPClient` | Membuat komunikasi HTTP dengan server Go. |
| `http.GET()` | Mengambil data dari API. |
| `deserializeJson()` | Membaca data JSON dari server. |
| `myServo.write()` | Menggerakkan motor servo sesuai sudut. |
| `delay(2000)` | Memberi jeda 2 detik sebelum membaca kembali. |

## 6. Pengujian Terpadu

### A. Mengetahui IP Address Komputer

Buka Command Prompt Windows:

```bash
ipconfig
```

Cari bagian IPv4 Address pada adapter WiFi atau Ethernet yang digunakan.

Contoh:

```text
IPv4 Address: 192.168.1.50
```

Masukkan alamat tersebut pada program ESP32:

```cpp
const char* serverUrl =
    "http://192.168.1.50:50001/servo/status";
```

Jangan menggunakan `localhost` pada ESP32 karena `localhost` akan mengarah ke ESP32 itu sendiri.

Pastikan:
- Komputer dan ESP32 terhubung ke jaringan yang sama.
- Server Go masih berjalan.
- Windows Firewall mengizinkan koneksi masuk pada TCP port `50001`.
- Alamat IP komputer benar.

### B. Upload Program ESP32

1. Hubungkan ESP32 ke komputer.
2. Buka Arduino IDE.
3. Pilih board ESP32.
4. Pilih port yang sesuai.
5. Masukkan program ESP32.
6. Klik Upload.
7. Buka Serial Monitor dengan baud rate `115200`.

Jika berhasil, akan muncul:

```text
Menghubungkan WiFi...
WiFi berhasil terhubung
IP ESP32: 192.168.1.100
```

### C. Menguji Pergerakan Servo

Pastikan server Go berjalan dan ESP32 sudah terhubung.

Buka Postman, kemudian kirim:

```http
PUT http://localhost:50001/servo/update/180
```

Server akan memperbarui nilai sudut di PostgreSQL.

ESP32 kemudian membaca nilai terbaru melalui API dan menggerakkan servo menuju sudut 180 derajat.

Ulangi pengujian menggunakan:

```http
PUT http://localhost:50001/servo/update/0
```

```http
PUT http://localhost:50001/servo/update/45
```

```http
PUT http://localhost:50001/servo/update/90
```

### D. Alur Kerja Sistem

```text
Postman
   |
   | PUT sudut servo
   v
Go REST API
   |
   | Menyimpan data
   v
PostgreSQL
   |
   | Mengirim status servo
   v
Go REST API
   |
   | HTTP GET
   v
ESP32
   |
   | Mengatur sudut
   v
Motor Servo
```

**Hasil akhir:** Ketika nilai sudut servo diubah melalui Postman, data akan disimpan di PostgreSQL. ESP32 membaca perubahan tersebut melalui REST API setiap 2 detik, kemudian menggerakkan motor servo sesuai nilai sudut yang diterima.

---

**Catatan:** Program backend pada contoh ini belum memvalidasi input sudut 0–180. Karena itu, gunakan nilai dalam rentang tersebut saat pengujian. Program ESP32 membatasi sudut gerakan menggunakan `constrain()`, tetapi validasi tersebut tidak mengubah nilai yang tersimpan di database.
