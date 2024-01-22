package migration

import (
	"database/sql"
	"fmt"
)

// TokenMigrate digunakan untuk menjalankan migrasi tabel token.
func TokenMigrate(db *sql.DB) error {

	// SQL statement untuk membuat tabel token dengan kolom "revoke"
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tokens (
            id CHAR(36) NOT NULL PRIMARY KEY,
		user_id CHAR(36),
            token VARCHAR(255) NOT NULL,
            created_at TIMESTAMP,
            updated_at TIMESTAMP,
            expired_at TIMESTAMP,
            is_revoked TINYINT(1) DEFAULT 0,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )
    `

	// Menjalankan perintah SQL untuk membuat tabel
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("gagal melakukan migrasi: %v", err)
	}

	return err
}
