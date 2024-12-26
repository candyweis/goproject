// internal/database/models.go
package database

import "time"

// User представляет пользователя в системе
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Score представляет запись с очками пользователя для конкретной игры
type Score struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	GameType  string    `json:"game_type"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
