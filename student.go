package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// 学生数据结构
type Student struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

var db *sql.DB

func main() {
	// 连接数据库
	var err error
	// Open不会检验用户名和密码
	if db, err = sql.Open("mysql", "gouser:Mxd20051020@@tcp(127.0.0.1:3306)/student_service"); err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}
	// Ping会检验用户名和密码
	if err = db.Ping(); err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}
	fmt.Println("数据库连接成功")
	defer db.Close()

	// 创建学生表
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS students (
            id     INT AUTO_INCREMENT PRIMARY KEY   COMMENT '编号',
            name   VARCHAR(50)            			COMMENT '姓名',
            age    TINYINT UNSIGNED       			COMMENT '年龄',
            gender CHAR(1)               	 		COMMENT '性别'
        ) ENGINE=InnoDB
          DEFAULT CHARSET=utf8mb4
          COLLATE=utf8mb4_0900_ai_ci
          COMMENT='学生信息表';
    `)
    if err != nil {
        fmt.Printf("学生信息表创建失败: %v\n", err)
    }
    fmt.Println("学生信息表创建完成")

	// 初始化Gin路由
	r := gin.Default()

	// 注册路由
	rGroup:=r.Group("/students")
	{
		// 获取所有学生
		rGroup.GET("", ListStudents)
		// 创建学生
		rGroup.POST("", CreateStudent)
		// 获取单个学生
		rGroup.GET("/:id", GetStudent)
		// 更新学生
		rGroup.PUT("/:id", UpdateStudent)
		// 删除学生
		rGroup.DELETE("/:id", DeleteStudent)
	}

	// 启动服务
	fmt.Println("StudentService 正在运行 :8080")
	r.Run(":8080")
}

// 创建学生
func CreateStudent(c *gin.Context) {
	var newStudent Student
	if err := c.BindJSON(&newStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	result, err := db.Exec("INSERT INTO students (name, age, gender) VALUES (?, ?, ?)",
		newStudent.Name, newStudent.Age, newStudent.Gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	newStudent.ID = int(id)

	c.JSON(http.StatusCreated, newStudent)
}
// 删除学生
func DeleteStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的学生ID"})
		return
	}

	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("学生 %d 已删除", id)})
}
// 更新学生
func UpdateStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的学生ID"})
		return
	}

	var updateData Student
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查学生是否存在
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM students WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
		return
	}

	_, err = db.Exec("UPDATE students SET name = ?, age = ?, gender = ? WHERE id = ?",
		updateData.Name, updateData.Age, updateData.Gender, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updateData.ID = id
	c.JSON(http.StatusOK, updateData)
}
// 获取单个学生
func GetStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的学生ID"})
		return
	}

	var s Student
	err = db.QueryRow("SELECT id, name, age, gender FROM students WHERE id = ?", id).
		Scan(&s.ID, &s.Name, &s.Age, &s.Gender)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}
// 获取所有学生
func ListStudents(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, age, gender FROM students")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Gender); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		students = append(students, s)
	}

	c.JSON(http.StatusOK, students)
}
