package handlers

import "net/http"

// ChooseGamePage — выбор между "Snake" и "Shooter"
func ChooseGamePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "choose.html", nil)
}

// SnakePage — страница со Змейкой
func SnakePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "snake.html", nil)
}

// ShooterPage — страница с Шутером
func ShooterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "shooter.html", nil)
}
