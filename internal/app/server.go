// internal/app/server.go
package app

import (
	"log"
	"net/http"
	"path/filepath"

	"my-go-project/internal/database"
	"my-go-project/internal/handlers"
	"my-go-project/internal/session"

	"github.com/gorilla/mux"
)

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

	// Применяем middleware для сессий
	router.Use(session.SessionMiddleware)

	// Публичные маршруты
	router.HandleFunc("/", handlers.IndexPage).Methods("GET")
	router.HandleFunc("/about", handlers.AboutPage).Methods("GET")
	router.HandleFunc("/contact", handlers.ContactPage).Methods("GET")
	router.HandleFunc("/register", handlers.RegisterPage(db)).Methods("GET", "POST")
	router.HandleFunc("/login", handlers.LoginPage(db)).Methods("GET", "POST")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")
	router.HandleFunc("/profile", handlers.ProfilePage(db)).Methods("GET")

	// Защищённые роуты (игры)
	gameRouter := router.PathPrefix("/game").Subrouter()
	gameRouter.Use(session.AuthMiddleware) // Проверяем авторизацию
	gameRouter.HandleFunc("/choose", handlers.ChooseGamePage).Methods("GET")
	gameRouter.HandleFunc("/snake", handlers.SnakePage).Methods("GET")
	gameRouter.HandleFunc("/shooter", handlers.ShooterPage).Methods("GET")

	// Маршруты для очков
	router.HandleFunc("/score", handlers.SaveScoreHandler(db)).Methods("POST")
	router.HandleFunc("/scores", handlers.GetUserScoresHandler(db)).Methods("GET")

	// Логирование всех маршрутов
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, err := route.GetMethods()
		if err != nil {
			methods = []string{}
		}
		path, err := route.GetPathTemplate()
		if err != nil {
			path = "unknown"
		}
		log.Printf("Registered route: %v %v", methods, path)
		return nil
	})

	log.Println("Сервер запущен на :8080")
	return http.ListenAndServe(":8080", router)
}
