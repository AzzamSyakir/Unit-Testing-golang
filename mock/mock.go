// MockDB.go
package mocks

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

type DBHandler interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	// Tambahkan metode-metode lain yang mungkin Anda perlukan
}
