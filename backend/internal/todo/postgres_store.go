package todo

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"
	"time"
	"fmt"

)

type PostgresStore struct {
	DB *sql.DB
}

func (p *PostgresStore) Add(task Task) error {
	_, err := p.DB.Exec("INSERT INTO tasks(task, priority, status, due_date) VALUES($1,$2,$3,$4)",
		task.Task, task.Priority, task.Status, task.DueDate)
	return err
}

func (p *PostgresStore) Delete(id int64) error {
	_, err := p.DB.Exec("DELETE FROM tasks WHERE id=$1", id)
	return err
}

func (p *PostgresStore) GetAll() ([]Task, error) {
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



// в postgres_store.go
func (p *PostgresStore) UpdateStatus(id int64, status bool) error {
    _, err := p.DB.Exec("UPDATE tasks SET status = $1 WHERE id = $2", status, id)
    return err
}


func (p *PostgresStore) FilterTasks(tasks []Task, priority, statusStr, dateFilter string) []Task {
	res := []Task{}

    var statusFilter *bool
    if statusStr == "true" {
        val := true
        statusFilter = &val
    } else if statusStr == "false" {
        val := false
        statusFilter = &val
    }

    now := time.Now()

    for _, t := range tasks {
		// priority
		p := strings.TrimSpace(priority)
		tp := strings.TrimSpace(t.Priority)
		fmt.Printf("1, Task: %s | priority filter: %s | match: %v\n", t.Priority, priority, strings.EqualFold(strings.TrimSpace(t.Priority), strings.TrimSpace(priority)))
		if (p == "" || p == "all" || strings.EqualFold(tp, p)) &&
		(statusFilter == nil || t.Status == *statusFilter) {
			res = append(res, t)
			fmt.Printf("2 ,Task: %s | priority filter: %s | match: %v\n", t.Priority, priority, strings.EqualFold(strings.TrimSpace(t.Priority), strings.TrimSpace(priority)))
		}

		
        // по статусу
        if statusStr != "" && statusStr != "all" {
            if statusStr == "active" && t.Status { // true = выполнено
                continue
            }
            if statusStr == "completed" && !t.Status { // false = активное
                continue
            }
        }
		

        // по дате
        switch dateFilter {
        case "today":
            if t.DueDate.Format("2006-01-02") != now.Format("2006-01-02") {
                continue
            }
        case "week":
            y1, w1 := now.ISOWeek()
            y2, w2 := t.DueDate.ISOWeek()
            if y1 != y2 || w1 != w2 {
                continue
            }
        case "overdue":
            if !t.DueDate.Before(now) {
                continue
            }
        }
        res = append(res, t)
    }

    return res
}

