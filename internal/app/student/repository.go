// package student

// import "database/sql"

// type Repository interface {
// 	Create(student *Student) error
// 	GetByID(id int) (*Student, error)
// 	GetAll() ([]Student, error)
// 	Update(id int, student *Student) error
// 	Delete(id int) error
// }

// type repository struct {
// 	db *sql.DB
// }

// func NewRepository(db *sql.DB) Repository {
// 	return &repository{db: db}//创建一个repository实例，并将传入的指向sql.DB的指针赋值给repository的db字段，返回实例的地址实际是指向repository的一个指针
// 							  //由于 repository 实现了 Repository 接口的所有方法，因此 *repository 可以被赋值给 Repository 类型的变量。
// }

// func (r *repository) Create(s *Student) error {
// 	res, err := r.db.Exec(
// 		"INSERT INTO students (name, age, gender) VALUES (?, ?, ?)",
// 		s.Name, s.Age, s.Gender,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	id, err := res.LastInsertId()
// 	if err != nil {
// 		return err
// 	}
// 	s.ID = int(id)
// 	return nil
// }

// func (r *repository) GetByID(id int) (*Student, error) {
// 	s := &Student{}
// 	err := r.db.QueryRow(
// 		"SELECT id, name, age, gender FROM students WHERE id = ?",
// 		id,
// 	).Scan(&s.ID, &s.Name, &s.Age, &s.Gender)
	
// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	return s, err
// }

// func (r *repository) GetAll() ([]Student, error) {
// 	rows, err := r.db.Query("SELECT id, name, age, gender FROM students")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var students []Student
// 	for rows.Next() {
// 		var s Student
// 		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Gender); err != nil {
// 			return nil, err
// 		}
// 		students = append(students, s)
// 	}
// 	return students, nil
// }

// func (r *repository) Update(id int, s *Student) error {
// 	_, err := r.db.Exec(
// 		"UPDATE students SET name=?, age=?, gender=? WHERE id=?",
// 		s.Name, s.Age, s.Gender, id,
// 	)
// 	return err
// }

// func (r *repository) Delete(id int) error {
// 	_, err := r.db.Exec("DELETE FROM students WHERE id=?", id)
// 	return err
// }



package student

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Repository interface {
	Create(student *Student) error
	GetByID(id uint) (*Student, error)
	GetAll() ([]Student, error)
	Update(id uint, student *Student) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create 创建新学生记录
func (r *repository) Create(s *Student) error {
	// 使用 GORM 的 Create 方法插入记录
	result := r.db.Create(s)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID 根据ID获取学生记录
func (r *repository) GetByID(id uint) (*Student, error) {
	var s Student
	// 使用 First 方法查询，自动处理 RecordNotFound 错误
	result := r.db.First(&s, id)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &s, nil
}

// GetAll 获取所有学生记录
func (r *repository) GetAll() ([]Student, error) {
	var students []Student
	// 使用 Find 方法获取所有记录
	result := r.db.Find(&students)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	// 处理空结果集
	if len(students) == 0 {
		return []Student{}, nil
	}
	
	return students, nil
}

// Update 更新学生记录
func (r *repository) Update(id uint, student *Student) error {
	// 1. 先检查记录是否存在
	var existing Student
	if err := r.db.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	
	// 2. 使用 Updates 方法部分更新字段，避免覆盖零值
	// 注意：Updates 方法只更新非零值字段
	result := r.db.Model(&Student{}).Where("id = ?", id).Updates(student)
	if result.Error != nil {
		return result.Error
	}
	
	// 3. 检查是否实际更新了记录
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	
	return nil
}

// Delete 删除学生记录
func (r *repository) Delete(id uint) error {
	// 使用 Delete 方法，添加 Where 条件确保安全
	result := r.db.Where("id = ?", id).Delete(&Student{})
	
	if result.Error != nil {
		return result.Error
	}
	
	// 检查是否实际删除了记录
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	
	return nil
}