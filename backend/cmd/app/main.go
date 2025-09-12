package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"TodoApp/backend/internal/delivery/https"
	"TodoApp/backend/internal/repository/postgres"
	_ "github.com/lib/pq"
	"TodoApp/backend/internal/storage"
	"TodoApp/backend/internal/usecase"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dataSource := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	// Подключение к БД
	db, err := storage.NewDB(dataSource)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	storage.InitDB(db)

	// Инициализация сервиса и обработчиков
	store := &postgres.PostgresTaskStore{DB: db}
	https.Svc = usecase.NewService(store)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/tasks", https.TasksHandler)
	http.HandleFunc("/tasks/add", https.AddTaskHandler)
	http.HandleFunc("/tasks/delete", https.DeleteTaskHandler)
	http.HandleFunc("/tasks/update", https.UpdateTaskHandler)
	http.HandleFunc("/tasks/filter", https.FilterTasksHandler)

	fmt.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}