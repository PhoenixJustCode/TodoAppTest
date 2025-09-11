package todo

import (
	"database/sql"
	_ "github.com/lib/pq"

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


func (p *PostgresStore) FilterTasks(tasks []Task, day int16, priority, statusStr string) []Task {
    res := []Task{}

    var statusFilter *bool
    if statusStr == "true" {
        val := true
        statusFilter = &val
    } else if statusStr == "false" {
        val := false
        statusFilter = &val
    }

    for _, t := range tasks {
        // if (day == 0 || t.Days == day) && 
        if (priority == "" || priority == "all" || t.Priority == priority) &&
            (statusFilter == nil || t.Status == *statusFilter) {
            res = append(res, t)
        }
    }
    return res
}


