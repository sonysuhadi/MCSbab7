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
	password = "sony123"
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
