// internal/session/session.go
package session

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// Store является хранилищем сессий
var Store *sessions.CookieStore

func init() {
	// Жестко закодированный секретный ключ для сессий
	// Рекомендуется использовать достаточно длинный и случайный ключ
	const sessionKey = "3q2+796tvuXp/EZn6lRjv2tj3s+LcQoFyO07iJ2T6+M="

	if sessionKey == "" {
		log.Fatal("SESSION_KEY не установлен")
	}

	Store = sessions.NewCookieStore([]byte(sessionKey))
	// Установка опций безопасности
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 дней
		HttpOnly: true,
		// Secure:   true, // Включите, если используете HTTPS
	}
}

// SessionMiddleware обеспечивает наличие сессии для каждого запроса
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Инициализация сессии
		_, err := Store.Get(r, "my-session")
		if err != nil {
			log.Printf("Ошибка при получении сессии: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware проверяет, авторизован ли пользователь
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := Store.Get(r, "my-session")
		if err != nil {
			log.Printf("Ошибка при получении сессии: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		user, ok := sess.Values["user"].(string)
		if !ok || user == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Пользователь авторизован, продолжаем обработку
		next.ServeHTTP(w, r)
	})
}
