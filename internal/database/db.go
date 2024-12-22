package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectDB — подключение к БД
func ConnectDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=password dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to DB successfully!")
	return db, nil
}
