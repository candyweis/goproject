package main

import (
	"log"

	"my-go-project/internal/app"
)

func main() {
	if err := app.RunServer(); err != nil {
		log.Fatal(err)
	}
}
