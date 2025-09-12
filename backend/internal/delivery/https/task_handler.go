package https

import (
	"TodoApp/backend/internal/domain"
	. "TodoApp/backend/internal/usecase"
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

	if err := Svc.AddTask(t); err != nil {
		http.Error(w, "Failed to add", http.StatusInternalServerError)
		return
	}

	// Trim валидация
	t.Task = strings.TrimSpace(t.Task)
	if t.Task == "" || len(t.Task) > LIMIT_CHAR { // лимит на длину
		http.Error(w, "Invalid task text", http.StatusBadRequest)
		return
	}

	if t.DueDate.IsZero(){
		http.Error(w, "Filed to add ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

// POST /tasks/delete  (body JSON: {"id": 123})
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
	
	json.NewEncoder(w).Encode(map[string]string{"status":"deleted"})
}


// POST /tasks/update  (body JSON: {"id":123, "status": true})
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

	json.NewEncoder(w).Encode(map[string]string{"status":"updated"})
}



func FilterTasksHandler(w http.ResponseWriter, r *http.Request) {
    priority := r.URL.Query().Get("priority") // all | low | medium | high
    dateFilter := r.URL.Query().Get("date")   // all | today | week | overdue
    sortBy := r.URL.Query().Get("sort")       // date | priority
    sortOrder := r.URL.Query().Get("order")   // asc | desc
    statusFilter := r.URL.Query().Get("status") // all | active | completed

	
    tasks, err := Svc.GetAllTasks() // достаём все задачи из БД
    if err != nil {
        http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
        return
    }

    var filtered []domain.Task
    now := time.Now()

    for _, t := range tasks {
		if statusFilter != "" && statusFilter != "all" {
            if statusFilter == "active" && t.Status { // true = выполнено
                continue
            }
            if statusFilter == "completed" && !t.Status { // false = активное
                continue
            }
        }
		// 🔹 Фильтр по приоритету
        if priority != "all" && priority != "" &&
            strings.ToLower(t.Priority) != priority {
            continue
        }

        // 🔹 Фильтр по дате
        switch dateFilter {
        case "today":
            if t.DueDate.Format("2006-01-02") != now.Format("2006-01-02") {
                continue
            }
        case "week":
            year1, week1 := now.ISOWeek()
            year2, week2 := t.DueDate.ISOWeek()
            if year1 != year2 || week1 != week2 {
                continue
            }
        case "overdue":
            if !t.DueDate.Before(now) {
                continue
            }
        }
        filtered = append(filtered, t)
    }

    // 🔹 Сортировка
    if sortBy != "" {
        sort.Slice(filtered, func(i, j int) bool {
            switch sortBy {
            case "date":
                if sortOrder == "desc" {
                    return filtered[i].DueDate.After(filtered[j].DueDate)
                }
                return filtered[i].DueDate.Before(filtered[j].DueDate)
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
