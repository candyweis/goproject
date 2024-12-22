package app

import (
	"net/http"

	"my-go-project/internal/session"
)

// AuthMiddleware проверяет, залогинен ли пользователь
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := session.Store.Get(r, "my-session")

		user := sess.Values["user"]
		if user == nil || user == "" {
			// Не залогинен — перенаправляем на /login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SessionMiddleware — если нужны действия с сессиями до хендлеров
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Можно что-то делать тут
		next.ServeHTTP(w, r)
	})
}
