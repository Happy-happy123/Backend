// --------------版本 1--------------
// package student

// import "errors"

// var (
// 	ErrNotFound = errors.New("学生不存在")
// )

// type Service interface {
// 	CreateStudent(req *CreateRequest) (*Student, error)
// 	GetStudent(id int) (*Student, error)
// 	GetAllStudents() ([]Student, error)
// 	UpdateStudent(id int, req *UpdateRequest) (*Student, error)
// 	DeleteStudent(id int) error
// }

// type service struct {
// 	repo Repository
// }

// func NewService(repo Repository) Service {
// 	return &service{repo: repo}
// }

// func (s *service) CreateStudent(req *CreateRequest) (*Student, error) {
// 	student := &Student{
// 		Name:   req.Name,
// 		Age:    req.Age,
// 		Gender: req.Gender,
// 	}
// 	return student, s.repo.Create(student)
// }

// func (s *service) GetStudent(id int) (*Student, error) {
// 	student, err := s.repo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if student == nil {
// 		return nil, ErrNotFound
// 	}
// 	return student, nil
// }

// func (s *service) GetAllStudents() ([]Student, error) {
// 	return s.repo.GetAll()
// }

// func (s *service) UpdateStudent(id int, req *UpdateRequest) (*Student, error) {
// 	student, err := s.GetStudent(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 更新字段
// 	if req.Name != "" {
// 		student.Name = req.Name
// 	}
// 	if req.Age > 0 {
// 		student.Age = req.Age
// 	}
// 	if req.Gender != "" {
// 		student.Gender = req.Gender
// 	}

// 	return student, s.repo.Update(id, student)
// }

// func (s *service) DeleteStudent(id int) error {
// 	if _, err := s.GetStudent(id); err != nil {
// 		return err
// 	}
// 	return s.repo.Delete(id)
// }

// package student

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// var (
// 	ErrNotFound = errors.New("学生不存在")
// )

// type Service interface {
// 	CreateStudent(req *CreateRequest) (*Student, error)
// 	GetStudent(id uint) (*Student, error)
// 	GetAllStudents() ([]Student, error)
// 	UpdateStudent(id uint, req *UpdateRequest) (*Student, error)
// 	DeleteStudent(id uint) error
// }

// type service struct {
// 	repo  Repository
// 	redis *redis.Client
// 	ctx   context.Context
// }

// func NewService(repo Repository, redis *redis.Client) Service {
// 	return &service{
// 		repo:  repo,
// 		redis: redis,
// 		ctx:   context.Background(),
// 	}
// }

// func (s *service) CreateStudent(req *CreateRequest) (*Student, error) {
// 	// 创建学生对象
// 	newStudent := &Student{
// 		Name:   req.Name,
// 		Age:    req.Age,
// 		Gender: req.Gender,
// 	}

// 	// 保存到数据库
// 	if err := s.repo.Create(newStudent); err != nil {
// 		return nil, fmt.Errorf("创建学生失败: %w", err)
// 	}

// 	return newStudent, nil
// }

// func (s *service) GetStudent(id uint) (*Student, error) {
// 	// 1. 尝试从Redis获取
// 	cacheKey := fmt.Sprintf("student:%d", id)
// 	cached, err := s.redis.Get(s.ctx, cacheKey).Result()
	
// 	// 缓存命中
// 	if err == nil {
// 		var student Student
// 		if err := json.Unmarshal([]byte(cached), &student); err == nil {
// 			return &student, nil
// 		}
// 	}
	
// 	// 2. 缓存未命中，查询数据库
// 	student, err := s.repo.GetByID(id)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return nil, ErrNotFound
// 		}
// 		return nil, fmt.Errorf("获取学生失败: %w", err)
// 	}
	
// 	// 3. 写入缓存（TTL: 5分钟）
// 	studentJSON, _ := json.Marshal(student)
// 	if err := s.redis.Set(s.ctx, cacheKey, studentJSON, 5*time.Minute).Err(); err != nil {
// 		// 记录错误但不中断流程
// 		fmt.Printf("警告: 无法写入缓存: %v\n", err)
// 	}
	
// 	return student, nil
// }

// func (s *service) GetAllStudents() ([]Student, error) {
// 	// 对于获取所有学生，通常不缓存（数据量大、更新频繁）
// 	students, err := s.repo.GetAll()
// 	if err != nil {
// 		return nil, fmt.Errorf("获取所有学生失败: %w", err)
// 	}
	
// 	// 处理空结果集
// 	if len(students) == 0 {
// 		return []Student{}, nil
// 	}
	
// 	return students, nil
// }

// func (s *service) UpdateStudent(id uint, req *UpdateRequest) (*Student, error) {
// 	// 1. 获取现有学生记录
// 	student, err := s.repo.GetByID(id)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return nil, ErrNotFound
// 		}
// 		return nil, fmt.Errorf("获取学生失败: %w", err)
// 	}
	
// 	// 2. 更新字段（仅更新提供的字段）
// 	if req.Name != "" {
// 		student.Name = req.Name
// 	}
// 	if req.Age > 0 {
// 		student.Age = req.Age
// 	}
// 	if req.Gender != "" {
// 		student.Gender = req.Gender
// 	}
	
// 	// 3. 更新数据库
// 	if err := s.repo.Update(id, student); err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return nil, ErrNotFound
// 		}
// 		return nil, fmt.Errorf("更新学生失败: %w", err)
// 	}
	
// 	// 4. 清除缓存
// 	cacheKey := fmt.Sprintf("student:%d", id)
// 	if err := s.redis.Del(s.ctx, cacheKey).Err(); err != nil {
// 		// 记录错误但不中断流程
// 		fmt.Printf("警告: 无法清除缓存: %v\n", err)
// 	}
	
// 	return student, nil
// }

// func (s *service) DeleteStudent(id uint) error {
// 	// 1. 检查学生是否存在
// 	if _, err := s.repo.GetByID(id); err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return ErrNotFound
// 		}
// 		return fmt.Errorf("获取学生失败: %w", err)
// 	}
	
// 	// 2. 删除数据库记录
// 	if err := s.repo.Delete(id); err != nil {
// 		return fmt.Errorf("删除学生失败: %w", err)
// 	}
	
// 	// 3. 清除缓存
// 	cacheKey := fmt.Sprintf("student:%d", id)
// 	if err := s.redis.Del(s.ctx, cacheKey).Err(); err != nil {
// 		// 记录错误但不中断流程
// 		fmt.Printf("警告: 无法清除缓存: %v\n", err)
// 	}
	
// 	return nil
// }


// --------------版本 2--------------
package student

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = errors.New("学生不存在")
)

type Service interface {
	CreateStudent(req *CreateRequest) (*Student, error)
	GetStudent(id uint) (*Student, error)
	GetAllStudents() ([]Student, error)
	UpdateStudent(id uint, req *UpdateRequest) (*Student, error)
	DeleteStudent(id uint) error
}

type service struct {
	repo  Repository
	redis *redis.Client
	ctx   context.Context
}

func NewService(repo Repository, redis *redis.Client) Service {
	return &service{
		repo:  repo,
		redis: redis,
		ctx:   context.Background(),
	}
}

func (s *service) CreateStudent(req *CreateRequest) (*Student, error) {
	newStudent := &Student{
		Name:   req.Name,
		Age:    req.Age,
		Gender: req.Gender,
	}

	if err := s.repo.Create(newStudent); err != nil {
		return nil, fmt.Errorf("创建学生失败: %w", err)
	}

	// 写入缓存
	cacheKey := fmt.Sprintf("student:%d", newStudent.ID)
	studentJSON, _ := json.Marshal(newStudent)
	if err := s.redis.Set(s.ctx, cacheKey, studentJSON, 5*time.Minute).Err(); err != nil {
		fmt.Printf("警告: 创建学生后写入缓存失败: %v\n", err)
	}
	
	return newStudent, nil
}

func (s *service) GetStudent(id uint) (*Student, error) {
	cacheKey := fmt.Sprintf("student:%d", id)
	cached, err := s.redis.Get(s.ctx, cacheKey).Result()
	
	if err == nil {
		var student Student
		if err := json.Unmarshal([]byte(cached), &student); err == nil {
			return &student, nil
		}
	}
	
	student, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("获取学生失败: %w", err)
	}
	
	studentJSON, _ := json.Marshal(student)
	if err := s.redis.Set(s.ctx, cacheKey, studentJSON, 5*time.Minute).Err(); err != nil {
		fmt.Printf("警告: 无法写入缓存: %v\n", err)
	}
	
	return student, nil
}

func (s *service) GetAllStudents() ([]Student, error) {
	// 尝试从缓存获取所有学生
	cacheKey := "students:all"
	cached, err := s.redis.Get(s.ctx, cacheKey).Result()
	
	if err == nil {
		var students []Student
		if err := json.Unmarshal([]byte(cached), &students); err == nil {
			return students, nil
		}
	}
	
	students, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("获取所有学生失败: %w", err)
	}
	
	if len(students) == 0 {
		return []Student{}, nil
	}
	
	// 写入缓存
	studentsJSON, _ := json.Marshal(students)
	if err := s.redis.Set(s.ctx, cacheKey, studentsJSON, 1*time.Minute).Err(); err != nil {
		fmt.Printf("警告: 无法写入所有学生缓存: %v\n", err)
	}
	
	return students, nil
}

func (s *service) UpdateStudent(id uint, req *UpdateRequest) (*Student, error) {
	student, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("获取学生失败: %w", err)
	}
	
	if req.Name != "" {
		student.Name = req.Name
	}
	if req.Age > 0 {
		student.Age = req.Age
	}
	if req.Gender != "" {
		student.Gender = req.Gender
	}
	
	if err := s.repo.Update(id, student); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("更新学生失败: %w", err)
	}
	
	cacheKey := fmt.Sprintf("student:%d", id)
	if err := s.redis.Del(s.ctx, cacheKey).Err(); err != nil {
		fmt.Printf("警告: 无法清除缓存: %v\n", err)
	}
	
	// 清除所有学生列表缓存
	if err := s.redis.Del(s.ctx, "students:all").Err(); err != nil {
		fmt.Printf("警告: 无法清除所有学生缓存: %v\n", err)
	}
	
	return student, nil
}

func (s *service) DeleteStudent(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("获取学生失败: %w", err)
	}
	
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("删除学生失败: %w", err)
	}
	
	cacheKey := fmt.Sprintf("student:%d", id)
	if err := s.redis.Del(s.ctx, cacheKey).Err(); err != nil {
		fmt.Printf("警告: 无法清除缓存: %v\n", err)
	}
	
	// 清除所有学生列表缓存
	if err := s.redis.Del(s.ctx, "students:all").Err(); err != nil {
		fmt.Printf("警告: 无法清除所有学生缓存: %v\n", err)
	}
	
	return nil
}


// --------------版本 3（无Redis）--------------
// package student

// import (
// 	"errors"
// 	"fmt"
// )

// // 定义全局错误变量
// var (
// 	ErrNotFound = errors.New("学生不存在")
// )

// // Service 接口定义
// type Service interface {
// 	CreateStudent(req *CreateRequest) (*Student, error)
// 	GetStudent(id uint) (*Student, error)
// 	GetAllStudents() ([]Student, error)
// 	UpdateStudent(id uint, req *UpdateRequest) (*Student, error)
// 	DeleteStudent(id uint) error
// }

// type noRedisService struct {
// 	repo Repository
// }

// func NewServiceWithoutRedis(repo Repository) Service {
// 	return &noRedisService{repo: repo}
// }

// func (s *noRedisService) CreateStudent(req *CreateRequest) (*Student, error) {
// 	newStudent := &Student{
// 		Name:   req.Name,
// 		Age:    req.Age,
// 		Gender: req.Gender,
// 	}

// 	if err := s.repo.Create(newStudent); err != nil {
// 		return nil, fmt.Errorf("创建学生失败: %w", err)
// 	}
// 	return newStudent, nil
// }

// func (s *noRedisService) GetStudent(id uint) (*Student, error) {
// 	student, err := s.repo.GetByID(id)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return nil, ErrNotFound
// 		}
// 		return nil, fmt.Errorf("获取学生失败: %w", err)
// 	}
// 	return student, nil
// }

// func (s *noRedisService) GetAllStudents() ([]Student, error) {
// 	students, err := s.repo.GetAll()
// 	if err != nil {
// 		return nil, fmt.Errorf("获取所有学生失败: %w", err)
// 	}
// 	return students, nil
// }

// func (s *noRedisService) UpdateStudent(id uint, req *UpdateRequest) (*Student, error) {
// 	student, err := s.repo.GetByID(id)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return nil, ErrNotFound
// 		}
// 		return nil, fmt.Errorf("获取学生失败: %w", err)
// 	}
	
// 	if req.Name != "" {
// 		student.Name = req.Name
// 	}
// 	if req.Age > 0 {
// 		student.Age = req.Age
// 	}
// 	if req.Gender != "" {
// 		student.Gender = req.Gender
// 	}
	
// 	if err := s.repo.Update(id, student); err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return nil, ErrNotFound
// 		}
// 		return nil, fmt.Errorf("更新学生失败: %w", err)
// 	}
	
// 	return student, nil
// }

// func (s *noRedisService) DeleteStudent(id uint) error {
// 	if err := s.repo.Delete(id); err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			return ErrNotFound
// 		}
// 		return fmt.Errorf("删除学生失败: %w", err)
// 	}
// 	return nil
// }