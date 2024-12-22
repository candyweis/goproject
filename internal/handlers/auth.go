package handlers

import (
	"database/sql"
	"net/http"

	"my-go-project/internal/database"
	"my-go-project/internal/session"
)

// LogoutHandler очищает сессию и редиректит на главную
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.Store.Get(r, "my-session")
	sess.Values["user"] = ""
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

// RegisterPage: GET показывает форму, POST создаёт пользователя в БД
func RegisterPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			renderTemplate(w, "register.html", nil)

		case http.MethodPost:
			email := r.FormValue("email")
			password := r.FormValue("password")

			// Проверяем, есть ли такой email
			var exists bool
			err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
			if err != nil {
				http.Error(w, "DB error", http.StatusInternalServerError)
				return
			}
			if exists {
				renderTemplate(w, "register.html", "Пользователь с таким E-mail уже существует.")
				return
			}

			// Вставляем запись
			_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, password)
			if err != nil {
				http.Error(w, "DB error", http.StatusInternalServerError)
				return
			}
			// Редирект на логин
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

// LoginPage: GET — форма входа, POST — проверка логина и пароля
func LoginPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			renderTemplate(w, "login.html", nil)

		case http.MethodPost:
			email := r.FormValue("email")
			password := r.FormValue("password")

			var user database.User
			err := db.QueryRow(
				"SELECT id, email, password FROM users WHERE email=$1", email,
			).Scan(&user.ID, &user.Email, &user.Password)

			if err != nil {
				if err == sql.ErrNoRows {
					renderTemplate(w, "login.html", "Пользователь не найден")
				} else {
					http.Error(w, "DB error", http.StatusInternalServerError)
				}
				return
			}

			// Сравниваем пароли (в реальном проекте — bcrypt и т.п.)
			if user.Password != password {
				renderTemplate(w, "login.html", "Неверный пароль")
				return
			}

			// Если ок, кладём в сессию
			sess, _ := session.Store.Get(r, "my-session")
			sess.Values["user"] = user.Email
			sess.Save(r, w)

			http.Redirect(w, r, "/game/choose", http.StatusFound)
		}
	}
}
