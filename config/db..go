package config

import (
	"database/sql"
	"fmt"
	"os"
	"testing-golang/application/migration"

	_ "github.com/go-sql-driver/mysql" // Import driver MySQL
	"github.com/rs/zerolog/log"
)

var DB *sql.DB

// InitDB digunakan untuk menghubungkan ke database.
func InitDB() *sql.DB {

	//baca env nya
	sqlInfo := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Membuka koneksi ke database
	db, err := sql.Open("mysql", sqlInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("Gagal membuka koneksi database")
	}

	// Memeriksa koneksi ke database
	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Gagal melakukan ping ke database")
	}
	// Panggil fungsi migrate untuk inisialisasi migrasi database
	if err := migration.UserMigrate(db); err != nil {
		log.Error().Err(err).Msg("Gagal melakukan migrasi user")
	}

	if err := migration.TokenMigrate(db); err != nil {
		log.Error().Err(err).Msg("Gagal melakukan migrasi token")
	}

	log.Info().Msg("Terhubung ke database!, migration Success")

	DB = db

	return db
}
func InitDBTest() *sql.DB {

	//baca env nya
	sqlInfo := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME_test"),
	)

	// Membuka koneksi ke database
	db, err := sql.Open("mysql", sqlInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("Gagal membuka koneksi database")
	}

	// Memeriksa koneksi ke database
	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Gagal melakukan ping ke database")
	}

	// Panggil fungsi migrate untuk inisialisasi migrasi database
	if err := migration.UserMigrate(db); err != nil {
		log.Error().Err(err).Msg("Gagal melakukan migrasi user")
	}

	if err := migration.TokenMigrate(db); err != nil {
		log.Error().Err(err).Msg("Gagal melakukan migrasi token")
	}

	return db
}
