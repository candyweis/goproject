// internal/database/db.go
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// ConnectDB устанавливает соединение с базой данных PostgreSQL
func ConnectDB() (*sql.DB, error) {
	// Жестко закодированные параметры подключения
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "password"
		dbname   = "mydb"
	)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %v", err)
	}

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %v", err)
	}

	log.Println("Подключение к базе данных успешно!")
	return db, nil
}
