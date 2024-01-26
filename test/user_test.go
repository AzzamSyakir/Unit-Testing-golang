package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing-golang/application/controller"
	"testing-golang/application/repositories"
	"testing-golang/application/router"
	"testing-golang/application/service"
)

var loginToken string

func TestRegisterAPI(t *testing.T) {
	// Panggil SetupTest function
	TestSetup(t)

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
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
}
func TestLoginApi(t *testing.T) {
	// Panggil SetupTest function
	TestSetup(t)

	// Buat server test
	server := httptest.NewServer(router.Router(globalDB))
	defer server.Close()

	// Buat client HTTP
	client := &http.Client{}

	// Request dengan data register
	request, err := http.NewRequest(http.MethodPost, server.URL+"/uu", bytes.NewBufferString(`{"email": "asa@gmail.com", "password": "rahasia"}`))
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

	// Continue with other checks as needed

	// Periksa body respons
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
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
	// Buat mock untuk http.Request
	userID := "test-123"
	requestURL := fmt.Sprintf("/users/%s", userID)
	HttpRequest := httptest.NewRequest(http.MethodPut, requestURL, nil)
	// Set authorization header with the token
	HttpRequest.Header.Set("Content-Type", "application/json")
	HttpRequest.Header.Set("Accept", "application/json")

	// Buat mock untuk http.ResponseWriter
	recorder := httptest.NewRecorder()

	// Buat mock db
	userRepository := repositories.NewUserRepository(globalDB) // menggunakan globalDB yang dideklarasikan di setupTest
	userService := service.NewUserService(*userRepository)
	userController := controller.NewUserController(*userService)

	// Panggil fungsi controller
	userController.GetUserController(recorder, HttpRequest)

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
	userController.UpdateUserController(recorder, request)

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
