// test_setup.go

package test

import (
	"database/sql"
	"testing"
	"testing-golang/config"

	"github.com/joho/godotenv"
)

var globalDB *sql.DB

func TestSetup(t *testing.T) {
	envPath := "/var/www/html/testing-golang/.env" // Sesuaikan dengan path env Anda
	if err := godotenv.Load(envPath); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	db := config.InitDBTest() // Menginisialisasi database test
	globalDB = db
	// create connection to database
}
