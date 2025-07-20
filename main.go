package main

import (
	"StudentService/db"
	"StudentService/handlers"
	"github.com/gin-gonic/gin"
	"StudentService/services"



	_ "github.com/go-sql-driver/mysql"

)

func main() {
	// 初始化数据库
	db.InitDB()
	db.CreateStudentTable() 
	defer db.DB.Close()
	// 创建Gin实例
	r := gin.Default()
	// 创建服务层实例
	studentService := services.NewStudentService(db.DB)
	studentHandler := handlers.NewStudentHandler(studentService)
	// 路由配置
	studentGroup := r.Group("/student")
	{
		studentGroup.GET("", studentHandler.ListStudents)          // 获取所有学生
		studentGroup.POST("", studentHandler.CreateStudent)       // 创建学生
		studentGroup.GET("/:id", studentHandler.GetStudent)       // 获取单个学生
		studentGroup.PUT("/:id", studentHandler.UpdateStudent)    // 更新学生
		studentGroup.DELETE("/:id", studentHandler.DeleteStudent) // 删除学生
	}
	// 启动服务
	r.Run(":8080")
}