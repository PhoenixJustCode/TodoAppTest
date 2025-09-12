package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"TodoApp/backend/pkg/repository/postgres"
	"TodoApp/backend/pkg/storage"
	"TodoApp/backend/pkg/usecase"
	"TodoApp/backend/pkg/domain"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// TaskBackend будет заменой HTTP API для Wails
type TaskBackend struct {
	Svc *usecase.Service
}

func NewTaskBackend(svc *usecase.Service) *TaskBackend {
	return &TaskBackend{Svc: svc}
}

// Методы, публичные для фронтенда
func (t *TaskBackend) GetAllTasks() ([]domain.Task, error) {
    tasks, err := t.Svc.GetAllTasks()
    if err != nil {
        log.Fatal("❌ Error fetching tasks:", err)
    } else {
        fmt.Printf("📦 Got %d tasks from DB\n", len(tasks))
    }

    return tasks, err
}


func (t *TaskBackend) AddTask(task domain.Task) error {
	if task.DueDate == "" {
        task.DueDate = "0001-01-01 00:00:00"
    }
	return t.Svc.AddTask(task)
}

func (t *TaskBackend) DeleteTask(id int64) error {
	return t.Svc.DeleteTask(id)
}

func (t *TaskBackend) UpdateStatus(id int64, status bool) error {
	return t.Svc.UpdateStatus(id, status)
}

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	dataSource := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	db, err := storage.NewDB(dataSource)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
    defer db.Close()

	store := &postgres.PostgresTaskStore{DB: db}
	svc := usecase.NewService(store)
	taskBackend := NewTaskBackend(svc)

	err = wails.Run(&options.App{
		Title:  "TodoApp",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind: []interface{}{
			taskBackend, // биндим бэкенд к фронтенду
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
