// internal/handlers/snake.go
package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// SnakePage обрабатывает страницу игры Snake
func SnakePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/snake.html")
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
}
