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
