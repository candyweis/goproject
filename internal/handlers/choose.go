// internal/handlers/choose.go
package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// ChooseGamePage обрабатывает страницу выбора игры
func ChooseGamePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/choose.html")
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
