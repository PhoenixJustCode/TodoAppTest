package usecase

import(
	"TodoApp/backend/pkg/domain"
	"TodoApp/backend/pkg/repository/postgres"
)

type Service struct {
	Store *postgres.PostgresTaskStore
}

func NewService(store *postgres.PostgresTaskStore) *Service {
	return &Service{Store: store}
}

func (s *Service) AddTask(t domain.Task) error {
	return s.Store.Add(t)
}

func (s *Service) DeleteTask(id int64) error {
	return s.Store.Delete(id)
}

func (s *Service) GetAllTasks() ([]domain.Task, error) {
	return s.Store.GetAll()
}

func (s *Service) UpdateStatus(id int64, status bool) error {
    return s.Store.UpdateStatus(id, status)
}

func (s *Service) FilterTasks(tasks []domain.Task, priority, status, dateFilter string) []domain.Task {
    return s.Store.FilterTasks(tasks, priority, status, dateFilter)
}
