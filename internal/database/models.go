package database

// User — сущность таблицы "users"
type User struct {
	ID       int
	Email    string
	Password string
}
