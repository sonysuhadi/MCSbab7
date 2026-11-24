# MCS BAB 7 – Membangun REST API Servo dengan Golang dan PostgreSQL

Tutorial membangun REST API untuk mengelola sudut motor servo menggunakan Golang, Gin Framework, PostgreSQL, dan ESP32.

Pada praktikum ini, Golang digunakan sebagai backend API, PostgreSQL sebagai penyimpanan data, sedangkan ESP32 mengambil data dari API melalui WiFi untuk menggerakkan motor servo secara otomatis.

## Daftar Isi

1. Persiapan Tools
2. Instalasi Golang
3. Membuat Project
4. Struktur Folder
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
15. Pengujian API Menggunakan Postman
16. Integrasi ESP32 dan Motor Servo
17. Pengujian Terpadu dengan Motor Servo Fisik
18. Memeriksa Database
19. Troubleshooting

---

# 1. Persiapan Tools

Sebelum memulai praktikum, siapkan beberapa aplikasi berikut.

| Tools | Fungsi |
|---|---|
| Golang | Bahasa pemrograman backend |
| Visual Studio Code | Editor untuk menulis kode |
| PostgreSQL | Database untuk menyimpan status servo |
| pgAdmin 4 | Mengelola database PostgreSQL |
| Postman | Menguji REST API |
| Arduino IDE | Menulis dan mengunggah program ESP32 |
| ESP32 DevKit | Mikrokontroler |
| Motor Servo | Aktuator untuk mengatur sudut putaran |

# 2. Instalasi Golang

Unduh Golang melalui:

https://go.dev/dl/

Setelah instalasi selesai, buka CMD atau PowerShell.

Periksa instalasi dengan perintah:

```bash
go version
```

Jika berhasil, versi Golang akan ditampilkan pada terminal.

# 3. Membuat Project

Buat folder project dengan nama:

```text
bab7
```

Contoh lokasi folder:

```text
C:\Users\sony\Documents\app flutter\bab7
```

Buka folder tersebut menggunakan Visual Studio Code.

Langkah-langkah:

1. Buka Visual Studio Code.
2. Pilih File → Open Folder.
3. Pilih folder `bab7`.
4. Buka terminal menggunakan `Ctrl + ``.

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

| Folder/File | Fungsi |
|---|---|
| `controllers` | Mengatur request dan response API |
| `database` | Mengatur koneksi dan migrasi database |
| `sql_migrations` | Menyimpan perintah SQL |
| `entities` | Mendefinisikan struktur data servo |
| `repositories` | Menjalankan query SQL |
| `routers` | Mengatur endpoint API |
| `main.go` | Menjalankan aplikasi |

# 5. Inisialisasi Go Module

Buka terminal pada folder project, kemudian jalankan:

```bash
go mod init bab7
```

Perintah tersebut membuat file `go.mod` yang berfungsi sebagai identitas module Go dan tempat pengelolaan dependency project.

# 6. Instalasi Dependencies

Jalankan perintah berikut:

```bash
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/rubenv/sql-migrate
```

Fungsi setiap library:

- `gin-gonic/gin`: Framework untuk membuat REST API.
- `lib/pq`: Driver untuk menghubungkan Golang dengan PostgreSQL.
- `sql-migrate`: Library untuk menjalankan migrasi database.

# 7. Membuat Entity

Buka folder `entities`, kemudian buat file:

`servo_entity.go`

Isi kode:

```go
package entities

type Servo struct {
    Id          int `json:"id"`
    StatusServo int `json:"srv_status"`
}
```

**Penjelasan kode:**

- `package entities`: Menandakan file termasuk dalam package entities.
- `type Servo struct`: Membuat struktur data bernama Servo.
- `Id`: Menyimpan identitas data servo.
- `StatusServo`: Menyimpan nilai sudut atau status servo.
- `json:"id"` dan `json:"srv_status"`: Menentukan nama properti ketika data diubah menjadi JSON.

Struktur ini digunakan ketika aplikasi mengambil atau memperbarui data servo.

# 8. Membuat Database dan Migrasi

Pada tahap ini, kita membuat fungsi untuk menjalankan migrasi database secara otomatis.

## A. Membuat file database.go

Masuk ke folder `database`.

Buat file:

`database.go`

Isi kode lengkap:

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

**Penjelasan fungsi:**

| Kode | Fungsi |
|---|---|
| `database/sql` | Mengelola koneksi database |
| `embed` | Menyertakan file SQL ke dalam aplikasi |
| `go:embed` | Membaca file migrasi dari folder `sql_migrations` |
| `DbConnection` | Menyimpan koneksi database |
| `DBMigrations()` | Menjalankan migrasi database |
| `migrate.Up` | Menjalankan migrasi ke versi terbaru |

Fungsi `DBMigrations()` digunakan agar tabel database dapat dibuat secara otomatis ketika aplikasi dijalankan.

## B. Membuat file migrasi SQL

Masuk ke folder:

`database/sql_migrations/`

Buat file:

`1_initiate.sql`

Isi kode:

```sql
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE servo (
    id INTEGER PRIMARY KEY,
    status_servo INTEGER
);

-- +migrate StatementEnd
```

**Penjelasan:**

- `CREATE TABLE servo`: Membuat tabel bernama servo.
- `id INTEGER PRIMARY KEY`: Menyimpan ID servo sebagai primary key.
- `status_servo INTEGER`: Menyimpan nilai status servo.

Struktur tabel yang dihasilkan:

| Kolom | Tipe Data | Fungsi |
|---|---|---|
| id | INTEGER | Identitas data servo |
| status_servo | INTEGER | Menyimpan nilai sudut servo |

# 9. Membuat Repository

Masuk ke folder `repositories`.

Buat file:

`servo_repo.go`

Isi kode lengkap:

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

**Penjelasan fungsi:**

### A. InitProj()

```go
func InitProj(db *sql.DB) (err error)
```

Digunakan untuk membuat data awal servo.

```sql
INSERT INTO servo (id, status_servo)
VALUES (1, 0)
ON CONFLICT (id) DO NOTHING
```

Artinya, sistem membuat data servo dengan ID `1` dan status awal `0`. Jika data sudah ada, sistem tidak membuat duplikat.

### B. GetStatus()

```go
func GetStatus(db *sql.DB) (result []entities.Servo, err error)
```

Digunakan untuk mengambil data servo dari PostgreSQL.

- `db.Query()`: Menjalankan perintah SELECT.
- `rows.Next()`: Membaca data setiap baris.
- `rows.Scan()`: Memasukkan hasil query ke dalam variabel.
- `result`: Menyimpan hasil data servo.

### C. UpdateStatus()

```go
func UpdateStatus(db *sql.DB, status entities.Servo) (err error)
```

Digunakan untuk memperbarui status servo.

Query:

```sql
UPDATE servo SET status_servo = $1 WHERE id = 1
```

Nilai `$1` diisi menggunakan `status.StatusServo`.

# 10. Membuat Controller

Masuk ke folder `controllers`.

Buat file:

`servo_controller.go`

Isi kode lengkap:

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

**Penjelasan fungsi:**

| Fungsi | Kegunaan |
|---|---|
| `InitProj()` | Memanggil repository untuk membuat data awal |
| `GetStatus()` | Mengambil status servo dan mengirimkan JSON |
| `UpdateStatus()` | Mengambil nilai dari URL dan memperbarui database |

Beberapa bagian penting:

- `gin.Context`: Menangani informasi request dan response.
- `c.Param("srv_status")`: Mengambil nilai dari parameter URL.
- `strconv.Atoi()`: Mengubah teks menjadi integer.
- `c.JSON()`: Mengirimkan response dalam format JSON.
- `http.StatusOK`: Menandakan request berhasil dengan HTTP 200.

# 11. Membuat Router

Masuk ke folder `routers`.

Buat file:

`servo_router.go`

Isi kode lengkap:

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

**Penjelasan:**

- `gin.Default()`: Membuat server Gin dengan middleware bawaan.
- `POST`: Digunakan untuk inisialisasi data servo.
- `GET`: Digunakan untuk mengambil status servo.
- `PUT`: Digunakan untuk memperbarui status servo.
- `:srv_status`: Parameter URL untuk menentukan nilai servo.

Daftar endpoint:

| Method | Endpoint | Fungsi |
|---|---|---|
| POST | `/servo/init-proj` | Membuat data awal |
| GET | `/servo/status` | Mengambil status servo |
| PUT | `/servo/update/:srv_status` | Memperbarui status servo |

# 12. Membuat Main Program

Buka file:

`main.go`

Isi kode lengkap sesuai struktur konfigurasi project:

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
        host, port, user, password, dbName,
    )

    DB, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal(err)
    }
    defer DB.Close()

    err = DB.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("DB berhasil tersambung")

    err = database.DBMigrations(DB)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Server berjalan di port ", PORT)

    err = routers.StartServer().Run(PORT)
    if err != nil {
        log.Fatal(err)
    }
}
```

Ganti `PASSWORD_POSTGRESQL_KAMU` dengan password PostgreSQL lokal.

**Penjelasan fungsi:**

1. `sql.Open()`: Membuka koneksi ke PostgreSQL.
2. `DB.Ping()`: Memeriksa apakah database dapat dihubungi.
3. `database.DBMigrations()`: Menjalankan migrasi tabel.
4. `routers.StartServer()`: Memanggil konfigurasi endpoint.
5. `Run(PORT)`: Menjalankan server pada port `50001`.
6. `log.Fatal()`: Menampilkan error dan menghentikan aplikasi jika terjadi kegagalan.

Catatan keamanan: Jangan mengunggah password database asli ke repository publik. Gunakan password lokal pada konfigurasi pengembangan dan lindungi kredensial sebelum membagikan source code.

# 13. Konfigurasi PostgreSQL

Buka pgAdmin 4.

Buat database baru dengan nama:

```sql
CREATE DATABASE bab7;
```

Konfigurasi database:

| Konfigurasi | Nilai |
|---|---|
| Host | localhost |
| Port | 5432 |
| Username | postgres |
| Password | Password PostgreSQL lokal |
| Database | bab7 |

Pastikan PostgreSQL sudah berjalan dan password pada `main.go` sesuai.

# 14. Menjalankan Aplikasi

Buka terminal pada folder project.

Jalankan:

```bash
go mod tidy
```

Perintah tersebut digunakan untuk melengkapi dependency yang dibutuhkan.

Kemudian jalankan:

```bash
go run main.go
```

Jika berhasil, terminal akan menampilkan pesan:

```text
DB berhasil tersambung
Migrasi berhasil dijalankan: 1 migrasi
Server berjalan di port :50001

[GIN-debug] POST /servo/init-proj
[GIN-debug] GET /servo/status
[GIN-debug] PUT /servo/update/:srv_status

[GIN-debug] Listening and serving HTTP on :50001
```

Server API dapat diakses melalui:

```text
http://localhost:50001
```

Biarkan terminal tetap terbuka selama proses pengujian.

# 15. Pengujian API Menggunakan Postman

Sebelum ESP32 digunakan, pastikan REST API berjalan dengan baik.

Base URL:

```text
http://localhost:50001
```

## A. Inisialisasi Data Servo

Atur request Postman:

- Method: `POST`
- URL: `http://localhost:50001/servo/init-proj`

Klik **Send**.

Fungsi: Membuat data awal servo dengan ID `1` dan status `0`.

## B. Melihat Status Servo

Atur request:

- Method: `GET`
- URL: `http://localhost:50001/servo/status`

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

Fungsi: Mengambil nilai status servo yang tersimpan di PostgreSQL.

## C. Memperbarui Status Servo

Atur request:

- Method: `PUT`
- URL: `http://localhost:50001/servo/update/90`

Contoh response:

```json
{
    "StatusServo": 90
}
```

Fungsi: Mengubah nilai status servo menjadi `90`.

Nilai dapat diganti, misalnya:

```text
/servo/update/0
/servo/update/45
/servo/update/90
/servo/update/180
```

# 16. Integrasi ESP32 dan Motor Servo

Pada tahap ini, ESP32 akan mengambil data dari REST API Golang melalui WiFi. Data JSON yang diterima akan dibaca untuk menentukan sudut motor servo.

## A. Peralatan

- ESP32 DevKit.
- Motor servo.
- Kabel jumper.
- Sumber daya servo sesuai spesifikasi motor.

## B. Instalasi Library Arduino IDE

Buka Arduino IDE.

Pilih:

`Sketch → Include Library → Manage Libraries`

Cari dan instal:

1. `ESP32Servo` oleh Kevin Harrington.
2. `ArduinoJson` oleh Benoit Blanchon versi 6 atau 7.

Library tersebut digunakan untuk mengendalikan servo dan membaca response JSON dari API.

## C. Rangkaian ESP32 dan Servo

Hubungkan komponen sebagai berikut:

| Pin Servo | ESP32 |
|---|---|
| Signal (kuning/oranye) | GPIO 13 |
| VCC (merah) | Sumber daya sesuai spesifikasi servo |
| GND (cokelat/hitam) | GND |

Pastikan GND sumber daya servo terhubung dengan GND ESP32. Jangan memasok daya motor servo dari pin 3.3V ESP32.

## D. Mengetahui IP Komputer

ESP32 harus mengakses alamat IP komputer yang menjalankan server Golang.

Pada Windows, buka CMD:

```bash
ipconfig
```

Cari bagian `IPv4 Address` pada adapter WiFi yang sedang digunakan.

Contoh:

```text
192.168.1.50
```

Alamat IP tersebut akan digunakan dalam konfigurasi API ESP32.

Jangan menggunakan `localhost`, karena localhost pada ESP32 menunjuk ke ESP32 itu sendiri, bukan komputer server.

## E. Program Lengkap ESP32

Buka Arduino IDE dan buat sketch baru.

Salin kode lengkap berikut:

```cpp
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include <ESP32Servo.h>

// 1. Konfigurasi WiFi
const char* ssid = "NAMA_WIFI_ANDA";
const char* password = "PASSWORD_WIFI_ANDA";

// 2. Alamat API Golang
// Ganti IP sesuai alamat komputer server
const char* serverUrl =
  "http://192.168.1.50:50001/servo/status";

// 3. Konfigurasi Servo
Servo myServo;

const int servoPin = 13;
int lastAngle = -1;

void setup() {
    Serial.begin(115200);
    delay(1000);

    // Alokasi timer PWM ESP32
    ESP32PWM::allocateTimer(0);
    ESP32PWM::allocateTimer(1);
    ESP32PWM::allocateTimer(2);
    ESP32PWM::allocateTimer(3);

    // Pengaturan servo
    myServo.setPeriodHertz(50);
    myServo.attach(servoPin, 500, 2400);

    // Menghubungkan ESP32 ke WiFi
    Serial.print("Menghubungkan ke WiFi: ");
    Serial.println(ssid);

    WiFi.begin(ssid, password);

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    Serial.println("\nWiFi berhasil terhubung!");

    Serial.print("IP ESP32: ");
    Serial.println(WiFi.localIP());
}

void loop() {

    // Memastikan WiFi terhubung
    if (WiFi.status() == WL_CONNECTED) {

        HTTPClient http;

        Serial.print("Menghubungi Server: ");
        Serial.println(serverUrl);

        http.begin(serverUrl);

        // Mengirim HTTP GET
        int httpResponseCode = http.GET();

        if (httpResponseCode == HTTP_CODE_OK) {

            String payload = http.getString();

            Serial.println("Response dari Server:");
            Serial.println(payload);

            // Parsing JSON
            DynamicJsonDocument doc(1024);

            DeserializationError error =
                deserializeJson(doc, payload);

            if (!error) {

                // Memeriksa data servo
                if (!doc["result"].is<JsonArray>() ||
                    doc["result"].size() == 0) {

                    Serial.println("Data servo tidak ditemukan.");

                } else {

                    // Mengambil nilai sudut servo
                    int targetAngle =
                        doc["result"][0]["srv_status"] | 0;

                    // Membatasi sudut 0 sampai 180
                    targetAngle = constrain(targetAngle, 0, 180);

                    // Menggerakkan servo jika nilai berubah
                    if (targetAngle != lastAngle) {

                        Serial.printf(
                            "Menggerakkan Servo ke %d derajat\n",
                            targetAngle
                        );

                        myServo.write(targetAngle);

                        lastAngle = targetAngle;

                    } else {

                        Serial.println(
                            "Sudut servo tidak berubah."
                        );
                    }
                }

            } else {

                Serial.print("Gagal parsing JSON: ");
                Serial.println(error.c_str());
            }

        } else {

            Serial.print("HTTP GET gagal, kode: ");
            Serial.println(httpResponseCode);
        }

        // Menutup koneksi HTTP
        http.end();

    } else {

        Serial.println("Koneksi WiFi terputus!");
        WiFi.reconnect();
    }

    // Mengambil data setiap 2 detik
    delay(2000);
}
```

Sebelum upload, ubah konfigurasi berikut:

```cpp
const char* ssid = "NAMA_WIFI_ANDA";
const char* password = "PASSWORD_WIFI_ANDA";
```

Isi dengan nama dan password WiFi yang digunakan.

Kemudian sesuaikan IP server:

```cpp
const char* serverUrl =
  "http://192.168.1.50:50001/servo/status";
```

Ganti `192.168.1.50` dengan IP komputer hasil perintah `ipconfig`.

**Catatan ArduinoJson:** Kode menggunakan `DynamicJsonDocument`, yang sesuai untuk ArduinoJson versi 6 dan masih tersedia pada versi 7 sebagai API kompatibilitas.

## F. Penjelasan Fungsi Program ESP32

| Bagian | Fungsi |
|---|---|
| `WiFi.h` | Menghubungkan ESP32 ke WiFi |
| `HTTPClient.h` | Mengirim request HTTP ke server |
| `ArduinoJson.h` | Membaca response JSON |
| `ESP32Servo.h` | Mengendalikan motor servo |
| `WiFi.begin()` | Memulai koneksi WiFi |
| `http.GET()` | Mengambil data dari API |
| `deserializeJson()` | Mengubah teks JSON menjadi data yang dapat dibaca |
| `doc["result"][0]["srv_status"]` | Mengambil nilai status servo |
| `constrain()` | Membatasi nilai sudut antara 0 dan 180 |
| `myServo.write()` | Menggerakkan servo |
| `delay(2000)` | Memberi jeda pengambilan data selama 2 detik |

## G. Cara Kerja Program

1. ESP32 dinyalakan.
2. ESP32 menghubungkan diri ke jaringan WiFi.
3. ESP32 mengirim request GET ke endpoint `/servo/status`.
4. Server Golang menerima request.
5. Repository mengambil data dari PostgreSQL.
6. Server mengirimkan data dalam format JSON.
7. ESP32 membaca nilai `srv_status`.
8. Servo bergerak sesuai nilai sudut yang diterima.
9. Proses diulangi setiap 2 detik.

# 17. Pengujian Terpadu dengan Motor Servo Fisik

Tahap ini digunakan untuk memastikan seluruh sistem bekerja mulai dari Postman, REST API, PostgreSQL, ESP32, hingga motor servo fisik.

Ikuti langkah berikut secara berurutan.

### Langkah 1: Mengetahui IP Komputer

Buka CMD pada komputer server, lalu jalankan:

```bash
ipconfig
```

Cari alamat IPv4 pada adapter WiFi.

Contoh:

```text
192.168.1.50
```

### Langkah 2: Mengubah Alamat Server pada ESP32

Masukkan alamat IP komputer ke variabel `serverUrl` pada program Arduino:

```cpp
const char* serverUrl =
  "http://192.168.1.50:50001/servo/status";
```

Pastikan IP tersebut sesuai dengan IP komputer yang menjalankan server Golang.

### Langkah 3: Menjalankan Server Golang

Buka terminal project dan jalankan:

```bash
go run main.go
```

Pastikan server berhasil tersambung ke PostgreSQL dan berjalan pada port `50001`.

### Langkah 4: Mengunggah Program ESP32

1. Hubungkan ESP32 ke komputer menggunakan kabel USB.
2. Buka Arduino IDE.
3. Pilih board ESP32 yang sesuai.
4. Pilih port COM.
5. Klik tombol Upload.
6. Tunggu sampai proses upload selesai.

### Langkah 5: Membuka Serial Monitor

Buka Serial Monitor pada Arduino IDE.

Atur baud rate menjadi:

```text
115200
```

Jika berhasil, ESP32 akan menampilkan informasi koneksi WiFi dan alamat IP-nya.

### Langkah 6: Mengubah Sudut Servo melalui Postman

Buka Postman dan buat request:

- Method: `PUT`
- URL: `http://localhost:50001/servo/update/180`

Klik **Send**.

Server akan menyimpan nilai `180` ke PostgreSQL.

### Langkah 7: Mengamati Motor Servo

ESP32 akan mengambil data terbaru dari endpoint GET setiap 2 detik.

Jika nilai `srv_status` berubah menjadi `180`, motor servo akan bergerak menuju sudut 180 derajat.

Ulangi pengujian dengan nilai berikut:

| Request PUT | Target Sudut |
|---|---|
| `/servo/update/0` | 0° |
| `/servo/update/45` | 45° |
| `/servo/update/90` | 90° |
| `/servo/update/180` | 180° |

**Hasil yang diharapkan:** Ketika nilai sudut diubah melalui Postman, data tersimpan di PostgreSQL. ESP32 membaca data tersebut melalui API dan menggerakkan motor servo secara otomatis.

Catatan pengujian:

- Komputer dan ESP32 harus berada pada jaringan lokal yang sama.
- Server Golang harus tetap aktif.
- Firewall Windows harus mengizinkan koneksi masuk pada port `50001`.
- Gunakan catu daya servo yang sesuai dengan spesifikasi motor.
- Jika servo tidak bergerak, periksa koneksi GND, kabel signal, dan sumber daya.

# 18. Memeriksa Database

Untuk memastikan nilai servo tersimpan dengan benar, buka pgAdmin 4.

Pilih database `bab7`, kemudian buka Query Tool.

Jalankan:

```sql
SELECT * FROM servo;
```

Contoh hasil:

| id | status_servo |
|---|---|
| 1 | 180 |

Jika nilai `status_servo` berubah setelah request PUT, berarti proses update berhasil.

# 19. Troubleshooting

| Masalah | Penyebab dan Solusi |
|---|---|
| Database gagal terhubung | Periksa PostgreSQL, port, username, dan password |
| `relation servo does not exist` | Pastikan migrasi berhasil dijalankan |
| `statusservo does not exist` | Gunakan nama kolom `status_servo` |
| `connection refused` | Pastikan PostgreSQL atau server API berjalan |
| Port 50001 sudah digunakan | Periksa aplikasi lain yang memakai port tersebut |
| `DB_PASSWORD` atau password salah | Periksa konfigurasi koneksi database |
| ESP32 gagal terhubung WiFi | Periksa SSID dan password WiFi |
| ESP32 gagal mengakses API | Periksa IP server, firewall, dan jaringan |
| JSON gagal dibaca | Periksa format response API dan library ArduinoJson |
| Servo tidak bergerak | Periksa GPIO, GND, kabel signal, dan catu daya servo |

---

# Alur Kerja Keseluruhan Sistem

```text
       POSTMAN
          |
          | PUT /servo/update/180
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
          | Menyimpan status servo
          |
          v
      ESP32 (WiFi)
          |
          | GET /servo/status
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
          | Response JSON
          v
         ESP32
          |
          v
      MOTOR SERVO
```

## Repository GitHub

Source code project:

https://github.com/sonysuhadi/MCSbab7

**Developer:** Sony Kurniawan Suhadi

**Project:** MCS BAB 7 – REST API Servo

Built with Golang, Gin Framework, PostgreSQL, and ESP32.
