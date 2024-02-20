package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"testing-golang/config"
	userHttp "testing-golang/internal/delivery/http"
	"testing-golang/internal/delivery/http/middleware"
	"testing-golang/internal/repository"
	"testing-golang/internal/usecase"

	"github.com/gorilla/mux"
)

func Router(db *sql.DB) *mux.Router {
	// Initialize repositories
	userRepository := repository.NewUserRepository(db)

	// Initialize services
	userusecase := usecase.NewUserUseCase(*userRepository)

	// Initialize http
	userhttp := userHttp.NewUserController(userusecase)

	// Create a new router
	router := mux.NewRouter()

	// Protected route
	protectedroute := router.PathPrefix("/").Subrouter()
	protectedroute.Use(middleware.AuthMiddleware)

	// Authentication route
	router.HandleFunc("/users", userhttp.Register).Methods("POST")
	router.HandleFunc("/users/login", userhttp.Login).Methods("POST")
	router.HandleFunc("/users/logout", userhttp.Logout).Methods("POST")

	// User route
	protectedroute.HandleFunc("/users", userhttp.Fetch).Methods("GET")
	protectedroute.HandleFunc("/users/{id}", userhttp.Get).Methods("GET")
	protectedroute.HandleFunc("/users/{id}", userhttp.Update).Methods("PUT")
	protectedroute.HandleFunc("/users/{id}", userhttp.Delete).Methods("DELETE")

	return router
}

func RunServer() {
	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to db: %w", err))
	}
	defer db.Close() // Pastikan koneksi database ditutup

	router := Router(db)

	// Mulai server HTTP dengan router yang telah dikonfigurasi
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
