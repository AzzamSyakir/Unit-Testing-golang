package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing-golang/application/controller"
	"testing-golang/application/entities"
	"testing-golang/application/repositories"
	"testing-golang/application/service"
	"testing-golang/config"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// Buat mock untuk http.Request
	request := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(`{"id": "0658b09c-6fbf-4eff-8aea-3243f837b09a", "password": "rahasia", "name": "asa", "email": "asa@gmail.com"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Buat mock untuk http.ResponseWriter
	recorder := httptest.NewRecorder()

	// Buat mock db
	db := config.InitDB() // Menginisialisasi database
	userRepository := repositories.NewUserRepository(db)
	userService := service.NewUserService(*userRepository)
	userController := controller.NewUserController(*userService)

	// Panggil fungsi controller
	userController.CreateUserController(recorder, request)

	// Assert hasil
	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Assert body response
	var user entities.User
	err := json.NewDecoder(response.Body).Decode(&user)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
