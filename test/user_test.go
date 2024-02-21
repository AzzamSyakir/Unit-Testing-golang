package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing-golang/cache"
	"testing-golang/config"
	"testing-golang/internal/delivery/http/router"
	"testing-golang/migrate"

	"github.com/joho/godotenv"
)

var loginToken string
var globalDB *sql.DB

func TestSetup(t *testing.T) {
	envpath := "../.env"
	err := godotenv.Load(envpath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize Redis
	cache.InitRedis(envpath)
	defer func() {
		if err := cache.RedisClient.Close(); err != nil {
			log.Println("Error closing Redis:", err)
		}
	}()
	db, err := config.InitDBTest() // Menginisialisasi database test
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	migrate.MigrateDB(db) // migrate tabel to database
	globalDB = db
	if globalDB == nil {
		t.Errorf("database null")
	}
}
func TestRegisterAPI(t *testing.T) {
	// Buat server test
	server := httptest.NewServer(router.Router(globalDB))
	defer server.Close()
	// Buat client HTTP
	client := &http.Client{}
	// Request dengan data register
	request, err := http.NewRequest(http.MethodPost, server.URL+"/users", bytes.NewBufferString(`{"id": "test-123", "password": "rahasia", "name": "tes", "email": "tes@gmail.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	// Kirim request
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	// Periksa status code
	if response.StatusCode == http.StatusNotFound {
		t.Error("404 page not found")
		return // Hentikan tes jika status code 404
	}
	if response.StatusCode != http.StatusCreated {
		var result map[string]interface{}
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Error decoding response body: %v", err)
		}
		errorMessage, ok := result["message"].(string)
		if !ok {
			t.Fatalf("Expected status code 200, got %d", response.StatusCode)
		} else {
			t.Fatalf("Expected status code 200, got %d. Error message: %s", response.StatusCode, errorMessage)
		}
		return
	}
	// Continue with other checks as needed
	// Periksa body respons
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
	t.Log("tes berhasil")
}

func TestLoginApi(t *testing.T) {
	// Buat server test
	server := httptest.NewServer(router.Router(globalDB))
	defer server.Close()
	// Buat client HTTP
	client := &http.Client{}
	// Request dengan data register
	request, err := http.NewRequest(http.MethodPost, server.URL+"/users/login", bytes.NewBufferString(`{"email": "tes@gmail.com","password": "rahasia"}`))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	// Kirim request
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	// Periksa status code
	if response.StatusCode == http.StatusNotFound {
		t.Error("404 page not found")
		return // Hentikan tes jika status code 404
	}
	if response.StatusCode != http.StatusCreated {
		var result map[string]interface{}
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Error decoding response body: %v", err)
		}
		errorMessage, ok := result["message"].(string)
		if !ok {
			t.Fatalf("Expected status code 200, got %d", response.StatusCode)
		} else {
			t.Fatalf("Expected status code 200, got %d. Error message: %s", response.StatusCode, errorMessage)
		}
		return
	}
	// Periksa body respons
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
	// Periksa token
	token, ok := result["data"].(map[string]interface{})["token"].(string)
	if !ok {
		t.Fatalf("Token tidak ditemukan dalam respons")
		return
	}
	// Simpan token login ke variabel global
	loginToken = token // Simpan token secara langsung
	t.Log("tes berhasil")
	// Pastikan token tidak kosong
	if loginToken == "" {
		t.Errorf("Token is empty")
	}
}
