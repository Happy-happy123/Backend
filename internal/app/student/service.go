package student

import "errors"

var (
	ErrNotFound = errors.New("学生不存在")
)

type Service interface {
	CreateStudent(req *CreateRequest) (*Student, error)
	GetStudent(id int) (*Student, error)
	GetAllStudents() ([]Student, error)
	UpdateStudent(id int, req *UpdateRequest) (*Student, error)
	DeleteStudent(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateStudent(req *CreateRequest) (*Student, error) {
	student := &Student{
		Name:   req.Name,
		Age:    req.Age,
		Gender: req.Gender,
	}
	return student, s.repo.Create(student)
}

func (s *service) GetStudent(id int) (*Student, error) {
	student, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, ErrNotFound
	}
	return student, nil
}

func (s *service) GetAllStudents() ([]Student, error) {
	return s.repo.GetAll()
}

func (s *service) UpdateStudent(id int, req *UpdateRequest) (*Student, error) {
	student, err := s.GetStudent(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		student.Name = req.Name
	}
	if req.Age > 0 {
		student.Age = req.Age
	}
	if req.Gender != "" {
		student.Gender = req.Gender
	}

	return student, s.repo.Update(id, student)
}

func (s *service) DeleteStudent(id int) error {
	if _, err := s.GetStudent(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}