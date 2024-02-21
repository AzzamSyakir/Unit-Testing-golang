package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing-golang/cache"
	"testing-golang/config"
	"testing-golang/internal/delivery/http/router"
	"testing-golang/migrate"

	"github.com/joho/godotenv"
)

var loginToken string
var globalDB *sql.DB

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
	if response.StatusCode != http.StatusOK {
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

	// Tutup client Redis setelah menyimpan token
	defer cache.RedisClient.Close() // Pindahkan penutupan ke sini
}

func TestFetchUserApi(t *testing.T) {
	// calling test login to set token value
	TestLoginApi(t)
	// get token login
	token := loginToken
	if token == "" {
		t.Fatal("Token tidak tersedia")
	}
	// Buat server test
	server := httptest.NewServer(router.Router(globalDB))
	defer server.Close()
	// Buat client HTTP
	client := &http.Client{}
	// Request dengan data register
	request, err := http.NewRequest(http.MethodGet, server.URL+"/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token) // Set authorization header dengan token yang sudah di-generate pada tes login sebelumnya
	// Kirim request
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	// Periksa status code
	if response.StatusCode == http.StatusNotFound {
		t.Error("404 page not found")
		return // Hentikan tes jika status code 404
	}
	if response.StatusCode != http.StatusOK {
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
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
}
func TestGetUserApi(t *testing.T) {
	// calling test login to set token value
	TestLoginApi(t)
	// get token login
	token := loginToken
	if token == "" {
		t.Fatal("Token tidak tersedia")
	}
	// calling SetupTest function
	// Buat server test
	server := httptest.NewServer(router.Router(globalDB))
	defer server.Close()
	// Buat client HTTP
	client := &http.Client{}
	userID := "test-123"
	requestURL := fmt.Sprintf("/users/%s", userID)
	// Request dengan data register
	request, err := http.NewRequest(http.MethodGet, server.URL+requestURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token) // Set authorization header dengan token yang sudah di-generate pada tes login sebelumnya
	// Kirim request
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	// Periksa status code
	if response.StatusCode == http.StatusNotFound {
		t.Error("404 page not found")
		return // Hentikan tes jika status code 404
	}
	if response.StatusCode != http.StatusOK {
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
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
}
func TestUpdateUserApi(t *testing.T) {
	// calling test login to set token value
	TestLoginApi(t)
	// get token login
	token := loginToken
	if token == "" {
		t.Fatal("Token tidak tersedia")
	}
	// calling SetupTest function
	// Buat server test
	server := httptest.NewServer(router.Router(globalDB))
	defer server.Close()
	// Buat client HTTP
	client := &http.Client{}
	userID := "test-123"
	requestBody := bytes.NewBufferString(`{"name": "update success"}`) // Tambahkan body untuk update nama
	requestURL := fmt.Sprintf("/users/%s", userID)
	// Request dengan data register
	request, err := http.NewRequest(http.MethodPut, server.URL+requestURL, requestBody)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token) // Set authorization header dengan token yang sudah di-generate pada tes login sebelumnya
	// Kirim request
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	// Periksa status code
	if response.StatusCode == http.StatusNotFound {
		t.Error("404 page not found")
		return // Hentikan tes jika status code 404
	}
	if response.StatusCode != http.StatusOK {
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
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
}
func TestLogoutUserApi(t *testing.T) {
	// calling test login to set token value
	TestLoginApi(t)
	// get token login
	token := loginToken
	if token == "" {
		t.Fatal("Token tidak tersedia")
	}

	// calling SetupTest function

	// Buat server test
	server := httptest.NewServer(router.Router(globalDB))
	defer server.Close()

	// Buat client HTTP
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, server.URL+("/users/logout"), nil)
	if err != nil {
		fmt.Printf("Error creating request, error: %v\n", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token) // Set authorization header dengan token yang sudah di-generate pada tes login sebelumnya

	// Kirim request
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	// Periksa status code
	if response.StatusCode == http.StatusNotFound {
		t.Error("404 page not found")
		return // Hentikan tes jika status code 404
	}

	if response.StatusCode != http.StatusOK {
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
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
}
func TestDeleteUserApi(t *testing.T) {
	// calling test login to set token value
	TestLoginApi(t)
	// get token login
	token := loginToken
	if token == "" {
		t.Fatal("Token tidak tersedia")
	}
	// calling SetupTest function
	// Buat server test
	server := httptest.NewServer(router.Router(globalDB))
	defer server.Close()
	// Buat client HTTP
	client := &http.Client{}
	userID := "test-123"
	requestURL := fmt.Sprintf("/users/%s", userID)
	request, err := http.NewRequest(http.MethodDelete, server.URL+requestURL, nil)
	if err != nil {
		fmt.Printf("Error creating request, error: %v\n", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token) // Set authorization header dengan token yang sudah di-generate pada tes login sebelumnya
	// Kirim request
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	// Periksa status code
	if response.StatusCode == http.StatusNotFound {
		t.Error("404 page not found")
		return // Hentikan tes jika status code 404
	}
	if response.StatusCode != http.StatusOK {
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
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

}

func TestMain(m *testing.M) {
	// Inisialisasi .env
	envpath := "../.env"
	if err := godotenv.Load(envpath); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	// Inisialisasi Redis
	cache.InitRedis(envpath)

	// Inisialisasi database test
	db, err := config.InitDBTest()
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}
	defer db.Close()

	// Migrasi tabel ke database
	if err := migrate.MigrateDB(db); err != nil {
		panic("Error migrating database: " + err.Error())
	}

	// Menjalankan semua tes
	exitCode := m.Run()

	// Mengakhiri koneksi Redis
	if err := cache.RedisClient.Close(); err != nil {
		panic("Error closing Redis: " + err.Error())
	}

	// Keluar dengan exit code yang sesuai
	os.Exit(exitCode)
}
