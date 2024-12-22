package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	layout := filepath.Join("templates", "layout.html")
	page := filepath.Join("templates", tmpl)

	log.Printf("Trying to parse layout: %s and page: %s", layout, page)
	t, err := template.ParseFiles(layout, page)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Templates parsed successfully. Executing...")

	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Template executed successfully.")
}

// IndexPage — главная
func IndexPage(w http.ResponseWriter, r *http.Request) {
	log.Println("IndexPage handler triggered")
	renderTemplate(w, "index.html", nil)
}
