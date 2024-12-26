// internal/handlers/score.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"my-go-project/internal/database"
	"my-go-project/internal/session"
)

// ScoreRequest представляет структуру запроса для сохранения очков
type ScoreRequest struct {
	GameType string `json:"game_type"`
	Score    int    `json:"score"`
}

// SaveScoreHandler обрабатывает сохранение очков
func SaveScoreHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверка метода запроса
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получение сессии
		sess, err := session.Store.Get(r, "my-session")
		if err != nil {
			log.Printf("Ошибка при получении сессии: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userEmail, ok := sess.Values["user"].(string)
		if !ok || userEmail == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Получение пользователя из базы данных
		var user database.User
		err = db.QueryRow("SELECT id, email FROM users WHERE email=$1", userEmail).Scan(&user.ID, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				log.Printf("Ошибка при запросе пользователя: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// Декодирование JSON-запроса
		var scoreReq ScoreRequest
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&scoreReq)
		if err != nil {
			log.Printf("Ошибка при декодировании JSON: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Валидация данных
		if scoreReq.GameType != "snake" && scoreReq.GameType != "shooter" {
			http.Error(w, "Invalid game type", http.StatusBadRequest)
			return
		}

		if scoreReq.Score < 0 {
			http.Error(w, "Invalid score", http.StatusBadRequest)
			return
		}

		// Вставка записи в таблицу scores
		_, err = db.Exec(
			"INSERT INTO scores (user_id, game_type, score) VALUES ($1, $2, $3)",
			user.ID, scoreReq.GameType, scoreReq.Score,
		)
		if err != nil {
			log.Printf("Ошибка при вставке очков: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Ответ клиенту
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

// GetUserScoresHandler возвращает список лучших очков пользователя
func GetUserScoresHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверка метода запроса
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получение сессии
		sess, err := session.Store.Get(r, "my-session")
		if err != nil {
			log.Printf("Ошибка при получении сессии: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userEmail, ok := sess.Values["user"].(string)
		if !ok || userEmail == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Получение пользователя из базы данных
		var user database.User
		err = db.QueryRow("SELECT id, email FROM users WHERE email=$1", userEmail).Scan(&user.ID, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				log.Printf("Ошибка при запросе пользователя: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// Получение лучших очков пользователя
		rows, err := db.Query(`
            SELECT game_type, MAX(score) as max_score
            FROM scores
            WHERE user_id = $1
            GROUP BY game_type
        `, user.ID)
		if err != nil {
			log.Printf("Ошибка при запросе очков: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type UserScore struct {
			GameType string `json:"game_type"`
			MaxScore int    `json:"max_score"`
		}

		var userScores []UserScore
		for rows.Next() {
			var us UserScore
			err := rows.Scan(&us.GameType, &us.MaxScore)
			if err != nil {
				log.Printf("Ошибка при сканировании очков: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			userScores = append(userScores, us)
		}

		// Проверка на ошибки после итерации
		if err = rows.Err(); err != nil {
			log.Printf("Ошибка строки: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Ответ клиенту
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userScores)
	}
}
