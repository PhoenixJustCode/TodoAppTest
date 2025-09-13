package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"
	"time"
	"sort"

	. "TodoApp/backend/pkg/domain"
)

type PostgresTaskStore struct {
	DB *sql.DB
}

// Добавление задачи
func (p *PostgresTaskStore) Add(task Task) error {
	_, err := p.DB.Exec(
		"INSERT INTO tasks(task, priority, status, due_date) VALUES($1,$2,$3,$4)",
		task.Task, task.Priority, task.Status, task.DueDate,
	)
	return err
}

// Удаление задачи по ID
func (p *PostgresTaskStore) Delete(id int64) error {
	_, err := p.DB.Exec("DELETE FROM tasks WHERE id=$1", id)
	return err
}

// Получение всех задач
func (p *PostgresTaskStore) GetAll() ([]Task, error) {
	rows, err := p.DB.Query("SELECT id, task, priority, status, due_date FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Task, &t.Priority, &t.Status, &t.DueDate); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// Обновление статуса задачи
func (p *PostgresTaskStore) UpdateStatus(id int64, status bool) error {
	_, err := p.DB.Exec("UPDATE tasks SET status = $1 WHERE id = $2", status, id)
	return err
}

// Фильтрация задач
func (p *PostgresTaskStore) FilterTasks(tasks []Task, priority, statusStr, dateFilter string) []Task {
	res := []Task{}
	now := time.Now()

	for _, t := range tasks {
		// 🔹 Фильтр по приоритету
		if priority != "" && priority != "all" && !strings.EqualFold(strings.TrimSpace(t.Priority), strings.TrimSpace(priority)) {
			continue
		}

		// Фильтр по статусу
		if statusStr != "" && statusStr != "all" {
			if statusStr == "active" && t.Status {
				continue
			}
			if statusStr == "completed" && !t.Status {
				continue
			}
		}

		// Фильтр по дате
		if t.DueDate != "" {
			due, err := time.Parse(time.RFC3339, t.DueDate)
			if err == nil {
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
		}

		res = append(res, t)
	}

	return res
}

func (p *PostgresTaskStore) SortTasks(tasks []Task, sortBy, sortOrder string) []Task {
	if sortBy == "" {
		return tasks
	}

	sort.Slice(tasks, func(i, j int) bool {
		switch sortBy {
		case "date":
			var di, dj time.Time
			if tasks[i].DueDate != "" {
				di, _ = time.Parse(time.RFC3339, tasks[i].DueDate)
			}
			if tasks[j].DueDate != "" {
				dj, _ = time.Parse(time.RFC3339, tasks[j].DueDate)
			}
			if sortOrder == "desc" {
				return di.After(dj)
			}
			return di.Before(dj)
		case "priority":
			order := map[string]int{"low": 1, "medium": 2, "high": 3}
			pi := order[strings.ToLower(tasks[i].Priority)]
			pj := order[strings.ToLower(tasks[j].Priority)]
			if sortOrder == "desc" {
				return pi > pj
			}
			return pi < pj
		default:
			return false
		}
	})
	return tasks
}
