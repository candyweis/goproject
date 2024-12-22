package app

import (
	"log"
	"net/http"
	"path/filepath"

	"my-go-project/internal/database"
	"my-go-project/internal/handlers"

	"github.com/gorilla/mux"
)

// RunServer запускает HTTP-сервер
func RunServer() error {
	// Подключение к PostgreSQL
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Инициализируем роутер
	router := mux.NewRouter()

	// Подключаем статические файлы на /static/
	staticDir := http.Dir(filepath.Join("static"))
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(staticDir)),
	)

	// Применяем middleware для сессий (опционально)
	router.Use(SessionMiddleware)

	// Публичные маршруты
	router.HandleFunc("/", handlers.IndexPage).Methods("GET")
	router.HandleFunc("/register", handlers.RegisterPage(db)).Methods("GET", "POST")
	router.HandleFunc("/login", handlers.LoginPage(db)).Methods("GET", "POST")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")

	// Защищённые роуты (игры)
	gameRouter := router.PathPrefix("/game").Subrouter()
	gameRouter.Use(AuthMiddleware) // Проверяем авторизацию
	gameRouter.HandleFunc("/choose", handlers.ChooseGamePage).Methods("GET")
	gameRouter.HandleFunc("/snake", handlers.SnakePage).Methods("GET")
	gameRouter.HandleFunc("/shooter", handlers.ShooterPage).Methods("GET")

	log.Println("Server started on :8080")
	return http.ListenAndServe(":8080", router)
}
