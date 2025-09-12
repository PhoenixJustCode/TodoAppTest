package https

import (
	"TodoApp/backend/pkg/domain"
	. "TodoApp/backend/pkg/usecase"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"
)

var Svc *Service
const LIMIT_CHAR = 1000

// GET /tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	tasks, err := Svc.GetAllTasks()
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// POST /tasks/add 
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var t domain.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Trim и валидация
	t.Task = strings.TrimSpace(t.Task)
	if t.Task == "" || len(t.Task) > LIMIT_CHAR {
		http.Error(w, "Invalid task text", http.StatusBadRequest)
		return
	}

	// Валидация даты
	if t.DueDate != "" {
		if _, err := time.Parse(time.RFC3339, t.DueDate); err != nil {
			http.Error(w, "Invalid due_date format", http.StatusBadRequest)
			return
		}
	}

	if err := Svc.AddTask(t); err != nil {
		http.Error(w, "Failed to add", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

// POST /tasks/delete
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct{ ID int64 `json:"id"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := Svc.DeleteTask(req.ID); err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// POST /tasks/update
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID     int64 `json:"id"`
		Status bool  `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := Svc.UpdateStatus(req.ID, req.Status); err != nil {
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// GET /tasks/filter
func FilterTasksHandler(w http.ResponseWriter, r *http.Request) {
	priority := r.URL.Query().Get("priority")
	dateFilter := r.URL.Query().Get("date")
	sortBy := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")
	statusFilter := r.URL.Query().Get("status")

	tasks, err := Svc.GetAllTasks()
	if err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}

	var filtered []domain.Task
	now := time.Now()

	for _, t := range tasks {
		// 🔹 Фильтр по статусу
		if statusFilter != "" && statusFilter != "all" {
			if statusFilter == "active" && t.Status {
				continue
			}
			if statusFilter == "completed" && !t.Status {
				continue
			}
		}

		// 🔹 Фильтр по приоритету
		if priority != "all" && priority != "" &&
			strings.ToLower(t.Priority) != strings.ToLower(priority) {
			continue
		}

		// 🔹 Фильтр по дате
		if t.DueDate != "" {
			due, err := time.Parse(time.RFC3339, t.DueDate)
			if err != nil {
				continue
			}

			switch dateFilter {
			case "today":
				if due.Format("2006-01-02") != now.Format("2006-01-02") {
					continue
				}
			case "week":
				y1, w1 := now.ISOWeek()
				y2, w2 := due.ISOWeek()
				if y1 != y2 || w1 != w2 {
					continue
				}
			case "overdue":
				if !due.Before(now) {
					continue
				}
			}
		}

		filtered = append(filtered, t)
	}

	// 🔹 Сортировка
	if sortBy != "" {
		sort.Slice(filtered, func(i, j int) bool {
			switch sortBy {
			case "date":
				var di, dj time.Time
				if filtered[i].DueDate != "" {
					di, _ = time.Parse(time.RFC3339, filtered[i].DueDate)
				}
				if filtered[j].DueDate != "" {
					dj, _ = time.Parse(time.RFC3339, filtered[j].DueDate)
				}
				if sortOrder == "desc" {
					return di.After(dj)
				}
				return di.Before(dj)
			case "priority":
				order := map[string]int{"low": 1, "medium": 2, "high": 3}
				pi := order[strings.ToLower(filtered[i].Priority)]
				pj := order[strings.ToLower(filtered[j].Priority)]
				if sortOrder == "desc" {
					return pi > pj
				}
				return pi < pj
			}
			return false
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}
