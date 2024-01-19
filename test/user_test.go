package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing-golang/application/controller"
	"testing-golang/application/repositories"
	"testing-golang/application/service"
	"testing-golang/config"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var loginToken string

func TestRegister(t *testing.T) {
	// initialize env
	envPath := "/var/www/html/testing-golang/.env" //absolute path to env file
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Buat mock untuk http.Request
	request := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(`{"id": "0658b09c-6fbf-4eff-8aea-3243f837b09a", "password": "rahasia", "name": "asa", "email": "asa@gmail.com"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Buat mock untuk http.ResponseWriter
	recorder := httptest.NewRecorder()

	db := config.InitDBTest() // Menginisialisasi database test
	userRepository := repositories.NewUserRepository(db)
	userService := service.NewUserService(*userRepository)
	userController := controller.NewUserController(*userService)

	// Panggil fungsi controller
	userController.CreateUserController(recorder, request)

	response := recorder.Result() // Dapatkan respons

	// Periksa status code
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	// Periksa body respons
	body, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)

	// Ambil data dari body respons
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Errorf("Error unmarshaling body: %v", err)
		return
	}

	// Periksa message
	expectedMessage := "Success" // Sesuaikan dengan pesan yang diinginkan
	actualMessage := data["message"].(string)
	assert.Equal(t, expectedMessage, actualMessage)

	t.Log("Tes berhasil")
}
func TestLogin(t *testing.T) {
	// initialize env
	envPath := "/var/www/html/testing-golang/.env" // absolute path to env file
	if err := godotenv.Load(envPath); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Buat mock untuk http.Request
	request := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"email": "asa@gmail.com", "password": "rahasia"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Buat mock untuk http.ResponseWriter
	recorder := httptest.NewRecorder()

	db := config.InitDBTest() // Menginisialisasi database test
	userRepository := repositories.NewUserRepository(db)
	userService := service.NewUserService(*userRepository)
	userController := controller.NewUserController(*userService)
	// Panggil fungsi controller
	userController.LoginUser(recorder, request)
	response := recorder.Result() // Dapatkan respons

	// Periksa status code
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

	// Continue with other checks as needed

	// Periksa body respons
	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Periksa token dan message
	token, ok := result["data"].(map[string]interface{})["token"].(string)
	if !ok {
		t.Fatalf("Token tidak ditemukan dalam respons")
		return
	}

	message, ok := result["message"].(string)
	if !ok {
		t.Fatalf("Pesan tidak ditemukan dalam respons")
		return
	}

	if message != "Login berhasil" {
		t.Fatalf("Pesan tidak sesuai: %s", message)
		return
	}

	// Simpan token login ke variabel global (jika diperlukan)
	loginToken = token
}

// func FetchUser(testing.T) {
// 	// initialize env
// 	envPath := "/var/www/html/testing-golang/.env" //absolute path to env file
// 	if err := godotenv.Load(envPath); err != nil {
// 		log.Fatal("Error loading .env file:", err)
// 	}

// 	// Buat mock untuk http.Request
// 	request := httptest.NewRequest(http.MethodGet, "/api/users")
// 	request.Header.Set("Content-Type", "application/json")
// 	request.Header.Set("Accept", "application/json")

// 	// Buat mock untuk http.ResponseWriter
// 	recorder := httptest.NewRecorder()

// 	// Buat mock db
// 	db := config.InitDBTest() // Menginisialisasi database test
// 	userRepository := repositories.NewUserRepository(db)
// 	userService := service.NewUserService(*userRepository)
// 	userController := controller.NewUserController(*userService)

// 	// Panggil fungsi controller
// 	userController.FetchUserController(recorder, request)

// 	response := recorder.Result() // Dapatkan respons

// 	// Periksa status code
// 	assert.Equal(t, http.StatusCreated, response.StatusCode)

// 	// Periksa body respons
// 	body, err := ioutil.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	// Ambil data dari body respons
// 	var data map[string]interface{}
// 	if err := json.Unmarshal(body, &data); err != nil {
// 		t.Errorf("Error unmarshaling body: %v", err)
// 		return
// 	}

// 	// Periksa message
// 	expectedMessage := "Success" // Sesuaikan dengan pesan yang diinginkan
// 	actualMessage := data["message"].(string)
// 	assert.Equal(t, expectedMessage, actualMessage)

// 	t.Log("Tes berhasil")
// }
