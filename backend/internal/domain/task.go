package domain

import "time"

type Task struct {
    ID       int64     `json:"id"`
    Task     string    `json:"task"`     // сам текст задания
    Priority string    `json:"priority"` // LOW / MEDIUM / HIGH
    Status   bool      `json:"status"`   // выполнено или нет
    DueDate  time.Time `json:"due_date"` // дата и время задачи
}
