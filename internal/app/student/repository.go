package student

import "database/sql"

type Repository interface {
	Create(student *Student) error
	GetByID(id int) (*Student, error)
	GetAll() ([]Student, error)
	Update(id int, student *Student) error
	Delete(id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}//创建一个repository实例，并将传入的指向sql.DB的指针赋值给repository的db字段，返回实例的地址实际是指向repository的一个指针
							  //由于 repository 实现了 Repository 接口的所有方法，因此 *repository 可以被赋值给 Repository 类型的变量。
}

func (r *repository) Create(s *Student) error {
	res, err := r.db.Exec(
		"INSERT INTO students (name, age, gender) VALUES (?, ?, ?)",
		s.Name, s.Age, s.Gender,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = int(id)
	return nil
}

func (r *repository) GetByID(id int) (*Student, error) {
	s := &Student{}
	err := r.db.QueryRow(
		"SELECT id, name, age, gender FROM students WHERE id = ?",
		id,
	).Scan(&s.ID, &s.Name, &s.Age, &s.Gender)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *repository) GetAll() ([]Student, error) {
	rows, err := r.db.Query("SELECT id, name, age, gender FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Gender); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (r *repository) Update(id int, s *Student) error {
	_, err := r.db.Exec(
		"UPDATE students SET name=?, age=?, gender=? WHERE id=?",
		s.Name, s.Age, s.Gender, id,
	)
	return err
}

func (r *repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM students WHERE id=?", id)
	return err
}