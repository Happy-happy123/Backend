// --------------版本 1--------------
// package main

// import (
// 	"fmt"
// 	"log"
// 	"StudentService/internal/app/student"
// 	"StudentService/internal/pkg/config"
// 	"StudentService/internal/pkg/database"
	// "StudentService/internal/router"
// )

// func main() {
	// 加载配置
	// cfg := config.Load()//返回一个结构体

					// // 初始化数据库
					// db, err := database.NewMySQLClient(cfg.DatabaseURL)
					// if err != nil {
					// 	log.Fatalf("数据库连接失败: %v", err)
					// }
					// defer db.Close()
					// fmt.Println("数据库连接成功")
					// fmt.Println("学生表创建成功")

					// // 初始化应用层
					// studentRepo := student.NewRepository(db)
					// studentSvc := student.NewService(studentRepo)
					// studentHandler := student.NewHandler(studentSvc)//返回一个指向Handler的指针

  	// 初始化MySQL
  	// gormDB, err := database.NewMySQLClient(cfg.MySQLDSN)
  	// if err != nil {
    	// log.Fatalf("MySQL连接失败: %v", err)
  	// }
  	// defer func() {
    	// sqlDB, _ := gormDB.DB()
    	// sqlDB.Close()
  	// }()
	// fmt.Println("MySQL连接成功")
	// fmt.Println("学生表创建成功")

  	// 初始化Redis
  	// redisClient := database.NewRedisClient(cfg.RedisAddr, "", cfg.RedisDB)
  	// defer redisClient.Close()

  	// 初始化应用层
  	// studentRepo := student.NewRepository(gormDB)
  	// studentSvc := student.NewService(studentRepo, redisClient)
  	// studentHandler := student.NewHandler(studentSvc)

	// 设置路由
	// r := router.SetupRouter(studentHandler)

	// 启动服务
	// fmt.Printf("StudentService 正在运行 :%d\n", cfg.Port)
	// r.Run(fmt.Sprintf(":%d", cfg.Port))
// }

// --------------版本 2--------------
package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // 导入pprof包
	"StudentService/internal/app/student"
	"StudentService/internal/pkg/config"
	"StudentService/internal/pkg/database"
	"StudentService/internal/router"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化MySQL
	gormDB, err := database.NewMySQLClient(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("MySQL连接失败: %v", err)
	}
	defer func() {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}()
	fmt.Println("MySQL连接成功")
	fmt.Println("学生表创建成功")

	// 初始化Redis
	redisClient := database.NewRedisClient(cfg.RedisAddr, "", cfg.RedisDB)
	defer redisClient.Close()

	// 初始化应用层
	studentRepo := student.NewRepository(gormDB)
	studentSvc := student.NewService(studentRepo, redisClient)
	studentHandler := student.NewHandler(studentSvc)

	// 设置路由
	r := router.SetupRouter(studentHandler)

	// 启动pprof性能分析服务器
	go func() {
		pprofAddr := "localhost:6060"
		log.Printf("启动pprof性能分析服务器: http://%s/debug/pprof", pprofAddr)
		if err := http.ListenAndServe(pprofAddr, nil); err != nil {
			log.Printf("pprof服务器启动失败: %v", err)
		}
	}()

	// 启动主服务
	serverAddr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("StudentService 正在运行 :%d\n", cfg.Port)
	
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}



// --------------版本 3(无Redis)--------------
// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	_ "net/http/pprof"
// 	"StudentService/internal/app/student"
// 	"StudentService/internal/pkg/config"
// 	"StudentService/internal/router"
// 	"gorm.io/gorm"
//    	"gorm.io/driver/mysql"
// )

// func main() {
// 	// 加载配置
// 	cfg := config.Load()

// 	// 初始化MySQL（无连接池优化）
// 	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{
// 		SkipDefaultTransaction: true, // 禁用默认事务
// 	})
// 	if err != nil {
// 		log.Fatalf("MySQL连接失败: %v", err)
// 	}
// 	defer func() {
// 		sqlDB, _ := db.DB()
// 		sqlDB.Close()
// 	}()

// 	// 初始化应用层（不使用Redis）
// 	studentRepo := student.NewRepository(db)
// 	studentSvc := student.NewServiceWithoutRedis(studentRepo)
// 	studentHandler := student.NewHandler(studentSvc)

// 	// 设置路由
// 	r := router.SetupRouter(studentHandler)

// 	// 启动pprof性能分析服务器
// 	go func() {
// 		pprofAddr := "localhost:6060"
// 		log.Printf("启动pprof性能分析服务器: http://%s/debug/pprof", pprofAddr)
// 		if err := http.ListenAndServe(pprofAddr, nil); err != nil {
// 			log.Printf("pprof服务器启动失败: %v", err)
// 		}
// 	}()

// 	// 启动主服务
// 	serverAddr := fmt.Sprintf(":%d", cfg.Port)
// 	fmt.Printf("StudentService 正在运行 (无Redis) :%d\n", cfg.Port)
	
// 	if err := r.Run(serverAddr); err != nil {
// 		log.Fatalf("服务启动失败: %v", err)
// 	}
// }