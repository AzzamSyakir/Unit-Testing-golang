package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing-golang/application/controller"
	"testing-golang/application/repositories"
	"testing-golang/application/service"

	"github.com/stretchr/testify/assert"
)

var loginToken string

func TestRegister(t *testing.T) {
	// call SetupTest function
	TestSetup(t)
	// Buat mock untuk http.Request
	request := httptest.NewRequest("POST", "/users", bytes.NewBufferString(`{"id": "test-23", "password": "rahasia", "name": "asa", "email": "asa@gmail.com"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Buat mock untuk http.ResponseWriter
	recorder := httptest.NewRecorder()

	userRepository := repositories.NewUserRepository(globalDB) //using globalDB that declare in setupTest
	userService := service.NewUserService(*userRepository)
	userController := controller.NewUserController(*userService)

	// Panggil fungsi controller
	userController.CreateUserController(recorder, request)

	response := recorder.Result() // Dapatkan respons

	// Periksa status code
	if response.StatusCode != http.StatusCreated {
		var data map[string]interface{}
		err := json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			t.Fatalf("Error decoding response body: %v", err)
		}

		errorMessage, ok := data["message"].(string)
		if !ok {
			t.Fatalf("Expected status code %d, got %d", http.StatusCreated, response.StatusCode)
		} else {
			t.Fatalf("Expected status code %d, got %d. Error message: %s", http.StatusCreated, response.StatusCode, errorMessage)
		}
		t.FailNow() // Menghentikan eksekusi tes saat ada kesalahan
		return
	}

	// Periksa body respons
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Error membaca body: %v", err)
		t.FailNow() // Menghentikan eksekusi tes saat ada kesalahan
		return
	}

	// Ambil data dari body respons
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("Error unmarshaling body: %v", err)
		t.FailNow() // Menghentikan eksekusi tes saat ada kesalahan
		return
	}

	// Periksa message
	expectedMessage := "Success" // Sesuaikan dengan pesan yang diinginkan
	actualMessage := data["message"].(string)
	assert.Equal(t, expectedMessage, actualMessage)

	t.Log("Tes berhasil")
}
func TestLogin(t *testing.T) {
	// call SetupTest function
	TestSetup(t)
	// Buat mock untuk http.Request
	request := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"email": "asa@gmail.com", "password": "rahasia"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Buat mock untuk http.ResponseWriter
	recorder := httptest.NewRecorder()

	userRepository := repositories.NewUserRepository(globalDB) //using globalDB that declare in setupTest
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
	t.Log("tes berhasil")
}
func TestFetchUser(t *testing.T) {
	// call SetupTest function
	TestSetup(t)
	// Buat mock untuk http.Request
	request := httptest.NewRequest("GET", "/users", nil)

	// Set authorization header with the token
	request.Header.Set("Authorization", "Bearer "+loginToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Buat mock untuk http.ResponseWriter
	recorder := httptest.NewRecorder()

	// Buat mock db
	userRepository := repositories.NewUserRepository(globalDB) // menggunakan globalDB yang dideklarasikan di setupTest
	userService := service.NewUserService(*userRepository)
	userController := controller.NewUserController(*userService)

	// Panggil fungsi controller
	userController.FetchUserController(recorder, request)

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

		t.FailNow() // Menghentikan eksekusi tes saat ada kesalahan
		return
	}

	// Periksa body respons
	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
		t.FailNow() // Menghentikan eksekusi tes saat ada kesalahan
		return
	}
}
func TestGetUser(t *testing.T) {
	// Panggil SetupTest function
	TestSetup(t)
	// Buat mock http request
	// Buat mock untuk http.Request
	userID := "test-123"
	requestURL := fmt.Sprintf("/users/%s", userID)
	request := httptest.NewRequest("GET", requestURL, nil)

	// Set authorization header with the token
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Buat mock untuk http.ResponseWriter
	recorder := httptest.NewRecorder()

	// Buat mock db
	userRepository := repositories.NewUserRepository(globalDB)
	userService := service.NewUserService(*userRepository)
	userController := controller.NewUserController(*userService)

	// Panggil fungsi controller
	userController.GetUserController(recorder, request)
	response := recorder.Result() // Dapatkan respons

	// Periksa status code
	if response.StatusCode != http.StatusOK {
		var result map[string]interface{}
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Error decoding response body: %v", err)
		}
		errorMessage := result["message"].(string)
		t.Fatalf("Expected status code 200, got %d. Error message: %s", response.StatusCode, errorMessage)
		t.FailNow() // Menghentikan eksekusi tes saat ada kesalahan
		return
	}

	// Periksa body respons
	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
		t.FailNow() // Menghentikan eksekusi tes saat ada kesalahan
		return
	}
	t.Log("Tes berhasil")
}
func TestUpdateUser(t *testing.T) {
	// call SetupTest function
	TestSetup(t)

	// Buat mock untuk http.Request
	userID := "test-123"
	requestURL := fmt.Sprintf("/api/users/%s", userID)
	requestBody := bytes.NewBufferString(`{"name": "update success"}`) // Tambahkan body untuk update nama
	request := httptest.NewRequest(http.MethodPut, requestURL, requestBody)
	// Set authorization header with the token
	request.Header.Set("Authorization", "Bearer "+loginToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Buat mock untuk http.ResponseWriter
	recorder := httptest.NewRecorder()

	// Buat mock db
	userRepository := repositories.NewUserRepository(globalDB) // menggunakan globalDB yang dideklarasikan di setupTest
	userService := service.NewUserService(*userRepository)
	userController := controller.NewUserController(*userService)

	// Panggil fungsi controller
	userController.FetchUserController(recorder, request)

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

		t.FailNow() // Menghentikan eksekusi tes saat ada kesalahan
		return
	}

	// Periksa body respons
	var result map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
		t.FailNow() // Menghentikan eksekusi tes saat ada kesalahan
		return
	}

	t.Log("Tes berhasil")
}
