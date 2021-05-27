package service

var (
	// Svc global service var
	Svc *Service
)

// Service struct
type Service struct {
	todoSvc TodoService
}

// New init service
func New() (s *Service) {
	s = &Service{
		todoSvc: NewTodoService(),
	}

	return s
}

// TodoSvc return todo service
func (s *Service) TodoSvc() TodoService {
	return s.todoSvc
}
