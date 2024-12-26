// internal/handlers/profile.go
package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"my-go-project/internal/database"
	"my-go-project/internal/session"
)

// ProfileData представляет данные для шаблона профиля
type ProfileData struct {
	Email      string         `json:"email"`
	BestScores map[string]int `json:"best_scores"`
}

// ProfilePage обрабатывает страницу профиля
func ProfilePage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("ProfilePage handler called") // Логирование

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

		bestScores := map[string]int{
			"snake":   0,
			"shooter": 0,
		}

		for rows.Next() {
			var gameType string
			var maxScore int
			err := rows.Scan(&gameType, &maxScore)
			if err != nil {
				log.Printf("Ошибка при сканировании очков: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			bestScores[gameType] = maxScore
		}

		// Проверка на ошибки после итерации
		if err = rows.Err(); err != nil {
			log.Printf("Ошибка строки: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Подготовка данных для шаблона
		profileData := ProfileData{
			Email:      user.Email,
			BestScores: bestScores,
		}

		// Рендеринг шаблона
		tmpl, err := template.ParseFiles("templates/layout.html", "templates/profile.html")
		if err != nil {
			log.Printf("Ошибка при парсинге шаблонов: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "layout", profileData)
		if err != nil {
			log.Printf("Ошибка при выполнении шаблона: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
