package todo

type Service struct {
	Store *PostgresStore
}

func NewService(store *PostgresStore) *Service {
	return &Service{Store: store}
}

func (s *Service) AddTask(t Task) error {
	return s.Store.Add(t)
}

func (s *Service) DeleteTask(id int64) error {
	return s.Store.Delete(id)
}

func (s *Service) GetAllTasks() ([]Task, error) {
	return s.Store.GetAll()
}

func (s *Service) UpdateStatus(id int64, status bool) error {
    return s.Store.UpdateStatus(id, status)
}

func (s *Service) FilterTasks(tasks []Task, priority, status, dateFilter string) []Task {
    return s.Store.FilterTasks(tasks, priority, status, dateFilter)
}
