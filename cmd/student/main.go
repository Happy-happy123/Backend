package main

import (
	"fmt"
	"log"
	"StudentService/internal/app/student"
	"StudentService/internal/pkg/config"
	"StudentService/internal/pkg/database"
	"StudentService/internal/router"
)

func main() {
	// 加载配置
	cfg := config.Load()//返回一个结构体

	// 初始化数据库
	db, err := database.NewMySQLClient(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()
	fmt.Println("数据库连接成功")
	fmt.Println("学生表创建成功")

	// 初始化应用层
	studentRepo := student.NewRepository(db)
	studentSvc := student.NewService(studentRepo)
	studentHandler := student.NewHandler(studentSvc)//返回一个指向Handler的指针

	// 设置路由
	r := router.SetupRouter(studentHandler)

	// 启动服务
	fmt.Printf("StudentService 正在运行 :%d\n", cfg.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}