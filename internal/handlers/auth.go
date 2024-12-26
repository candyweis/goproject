// internal/handlers/auth.go
package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"my-go-project/internal/session"

	"golang.org/x/crypto/bcrypt"
)

// RegisterPage обрабатывает регистрацию пользователя
func RegisterPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("templates/layout.html", "templates/register.html")
			if err != nil {
				log.Printf("Ошибка при парсинге шаблонов: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
			err = tmpl.ExecuteTemplate(w, "layout", nil)
			if err != nil {
				log.Printf("Ошибка при выполнении шаблона: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			// Валидация данных
			if email == "" || password == "" {
				http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
				return
			}

			// Хэширование пароля
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Ошибка при хэшировании пароля: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}

			// Вставка пользователя в базу данных
			_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, string(hashedPassword))
			if err != nil {
				log.Printf("Ошибка при вставке пользователя: %v", err)
				http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
				return
			}

			// Установка сессии
			sess, err := session.Store.Get(r, "my-session")
			if err != nil {
				log.Printf("Ошибка при получении сессии: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
			sess.Values["user"] = email
			err = sess.Save(r, w)
			if err != nil {
				log.Printf("Ошибка при сохранении сессии: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}

			// Перенаправление на страницу профиля
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
	}
}

// LoginPage обрабатывает вход пользователя
func LoginPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("templates/layout.html", "templates/login.html")
			if err != nil {
				log.Printf("Ошибка при парсинге шаблонов: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
			err = tmpl.ExecuteTemplate(w, "layout", nil)
			if err != nil {
				log.Printf("Ошибка при выполнении шаблона: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			// Валидация данных
			if email == "" || password == "" {
				http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
				return
			}

			// Получение хэшированного пароля из базы данных
			var storedHashedPassword string
			err := db.QueryRow("SELECT password FROM users WHERE email=$1", email).Scan(&storedHashedPassword)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
				} else {
					log.Printf("Ошибка при запросе пользователя: %v", err)
					http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				}
				return
			}

			// Сравнение паролей
			err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
			if err != nil {
				http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
				return
			}

			// Установка сессии
			sess, err := session.Store.Get(r, "my-session")
			if err != nil {
				log.Printf("Ошибка при получении сессии: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
			sess.Values["user"] = email
			err = sess.Save(r, w)
			if err != nil {
				log.Printf("Ошибка при сохранении сессии: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}

			// Перенаправление на страницу профиля
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
	}
}

// LogoutHandler обрабатывает выход пользователя из системы
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := session.Store.Get(r, "my-session")
	if err != nil {
		log.Printf("Ошибка при получении сессии: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Удаление пользователя из сессии
	delete(sess.Values, "user")
	err = sess.Save(r, w)
	if err != nil {
		log.Printf("Ошибка при сохранении сессии: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Перенаправление на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
