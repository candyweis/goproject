// cmd/main.go
package main

import (
	"log"

	"my-go-project/internal/app"
)

func main() {
	err := app.RunServer()
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
