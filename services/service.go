package services

import (
	"database/sql"
	"StudentService/models"

	"fmt"
	"strings"
)

// 学生服务结构体
type StudentService struct {
	db *sql.DB
}

// NewStudentService 创建学生服务实例
func NewStudentService(db *sql.DB) *StudentService {
	return &StudentService{db: db} 
}

// CreateStudent 创建学生
func (s *StudentService) CreateStudent(student *models.Student) (int64, error) {
	result, err := s.db.Exec(
		"INSERT INTO student (name, age, gender) VALUES (?, ?, ?)", 
		student.Name, student.Age, student.Gender,
	)
	if err != nil {
		return 0, fmt.Errorf("创建学生失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取学生ID失败: %w", err)
	}

	fmt.Printf("成功创建学生，ID: %d", id)
	return id, nil
}

// GetStudent 获取单个学生
func (s *StudentService) GetStudent(id int) (*models.Student, error) {
	student := &models.Student{}
	err := s.db.QueryRow(
		"SELECT id, name, age, gender FROM student WHERE id = ?", 
		id,
	).Scan(&student.Id, &student.Name, &student.Age, &student.Gender)

	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("学生不存在 (ID: %d)", id)
	case err != nil:
		return nil, fmt.Errorf("查询学生失败: %w", err)
	}

	return student, nil
}

// ListStudents 获取所有学生
func (s *StudentService) ListStudents() ([]models.Student, error) {
	rows, err := s.db.Query("SELECT id, name, age, gender FROM student")
	if err != nil {
		return nil, fmt.Errorf("查询学生列表失败: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Age, &student.Gender); err != nil {
			return nil, fmt.Errorf("扫描学生数据失败: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历学生列表失败: %w", err)
	}

	return students, nil
}

// UpdateStudent 更新学生
func (s *StudentService) UpdateStudent(id int, update *models.StudentUpdate) error {
	query := "UPDATE student SET"
	var args []interface{}
	var updates []string

	if update.Name != nil {
		updates = append(updates, " name = ?")
		args = append(args, *update.Name)
	}
	if update.Age != nil {
		updates = append(updates, " age = ?")
		args = append(args, *update.Age)
	}
	if update.Gender != nil {
		updates = append(updates, " gender = ?")
		args = append(args, *update.Gender)
	}

	if len(updates) == 0 {
		return fmt.Errorf("没有提供更新字段")
	}

	query += strings.Join(updates, ",") + " WHERE id = ?"
	args = append(args, id)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("更新学生失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("学生不存在 (ID: %d)", id)
	}

	fmt.Printf("成功更新学生，ID: %d，影响行数: %d", id, rowsAffected)
	return nil
}

// DeleteStudent 删除学生
func (s *StudentService) DeleteStudent(id int) error {
	result, err := s.db.Exec("DELETE FROM student WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("删除学生失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("学生不存在 (ID: %d)", id)
	}

	fmt.Printf("成功删除学生，ID: %d", id)
	return nil
}