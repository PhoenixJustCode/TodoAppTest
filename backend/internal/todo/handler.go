package todo

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"
)

var Svc *Service

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := Svc.GetAllTasks()
	if err != nil {
		http.Error(w, "Failed to get tasks", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm error", 400)
		return
	}

	task := r.FormValue("task")
	priority := r.FormValue("priority")
	statusStr := r.FormValue("status")
	dueDateStr := r.FormValue("due_date") // ожидаем ISO-формат или "yyyy-mm-ddThh:mm"
	status := false
	if statusStr == "true" || statusStr == "on" {
		status = true
	}
	// Парсим строку в time.Time
	var dueDate time.Time
	if dueDateStr != "" {
		var err error
		dueDate, err = time.Parse("2006-01-02T15:04", dueDateStr) // формат для input type="datetime-local"
		if err != nil {
			http.Error(w, "Invalid due date format", 400)
			return
		}
	}

	t := Task{
		Task:     task,
		Priority: priority,
		Status:   status,
		DueDate:  dueDate,
	}


	err := Svc.AddTask(t)
	if err != nil {
		http.Error(w, "Failed to add task", 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm error", 400)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	if err := Svc.DeleteTask(id); err != nil {
		http.Error(w, "Failed to delete task", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
