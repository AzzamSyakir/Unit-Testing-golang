package main

import (
	"fmt"
	"log"
	"os"
	"testing-golang/cache"
	"testing-golang/config"
	"testing-golang/internal/delivery/http/router"
	"testing-golang/migrate"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize Redis
	cache.InitRedis()
	defer func() {
		if err := cache.RedisClient.Close(); err != nil {
			log.Println("Error closing Redis:", err)
		}
	}()

	// Cek argumen command line
	if len(os.Args) > 1 {
		// Jika argumen adalah "migrate", jalankan migrasi
		if os.Args[1] == "migrate" {
			// Initialize database
			db, err := config.InitDB()
			if err != nil {
				log.Fatal("Error connecting to database:", err)
			}
			err = migrate.MigrateDB(db)
			if err != nil {
				log.Fatal("Error running migrations:", err)
			}
			fmt.Println("Migrations successfully!")
			return
		}
	}

	// Jika tidak ada argumen, jalankan server
	fmt.Println("Server started on port 9000")
	router.RunServer()
}
