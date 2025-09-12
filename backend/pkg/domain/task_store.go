package domain

type TaskStore interface {
    Add(task Task) error
    Delete(id int) error
    GetAll() ([]Task, error)
}
