package main

import (
	"log"

	"todoapp/assets"        // твой assets.go
	"TodoApp/backend/internal/usecase" // структура App для биндинга

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func main() {
	myApp := usecase.NewApp() // создаём структуру для биндинга

	err := wails.Run(&options.App{
		Title:  "TodoApp",
		Width:  1000,
		Height: 700,
		Assets: assets.Assets,           // embed фронта
		Bind:   []interface{}{myApp},   // биндим методы к JS
	})

	if err != nil {
		log.Fatal(err)
	}
}
