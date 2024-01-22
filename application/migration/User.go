package migration

import (
	"database/sql"
	"fmt"
)

func UserMigrate(db *sql.DB) error {
	// SQL statement untuk membuat tabel users
	createTableSQL := `
	    CREATE TABLE IF NOT EXISTS users (
		  id CHAR(36) NOT NULL PRIMARY KEY,
		  name VARCHAR(255) NOT NULL,
		  email VARCHAR(255) NOT NULL UNIQUE,
		  password VARCHAR(255) NOT NULL,
		  created_at TIMESTAMP NOT NULL,
		  updated_at TIMESTAMP NOT NULL
	    )
	`

	// Menjalankan perintah SQL untuk membuat tabel
	_, err := db.Exec(createTableSQL)
	if err != nil {
		// Menangani kesalahan jika terjadi kesalahan saat migrasi
		return fmt.Errorf("gagal melakukan migrasi: %v", err)
	}

	return err
}
