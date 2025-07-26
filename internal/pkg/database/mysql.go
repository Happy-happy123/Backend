// package database

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/go-sql-driver/mysql"
// )
// var db *sql.DB
// var err error
// func NewMySQLClient(s string)(*sql.DB,error){
// 	if db,err=sql.Open("mysql",s);err!=nil{// Open不会检验用户名和密码
// 		return nil,fmt.Errorf("sql.open failed:%v",err)
// 	}
// 	if err = db.Ping(); err != nil {// Ping会检验用户名和密码
// 		return nil,fmt.Errorf("db.ping failed:%v",err)
// 	}
// 	if _, err := db.Exec(`
//         CREATE TABLE IF NOT EXISTS students (
//             id     INT AUTO_INCREMENT PRIMARY KEY   COMMENT '编号',
//             name   VARCHAR(50)            			COMMENT '姓名',
//             age    TINYINT UNSIGNED       			COMMENT '年龄',
//             gender CHAR(1)               	 		COMMENT '性别'
//         ) ENGINE=InnoDB
//           DEFAULT CHARSET=utf8mb4
//           COLLATE=utf8mb4_0900_ai_ci
//           COMMENT='学生信息表';
//     `);err!=nil{
// 		return nil,fmt.Errorf("学生表创建失败:%v",err)
// 	}
// 	return db,nil
// }

package database

import (
	"fmt"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
type Student struct {//database 包不应该依赖 app/student 包,迁移只需要表结构，数据库层不需要知道业务层的完整模型定义
	ID     uint   `gorm:"primarykey"`
	Name   string
	Age    int
	Gender string
}
//gorm与mysql数据库连接
func NewMySQLClient(dsn string) (*gorm.DB, error) {
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
  })
  if err != nil {
    return nil, err
  }

//自动迁移表结构
	if err := db.AutoMigrate(&Student{}); err != nil {
		return nil, fmt.Errorf("学生信息表迁移失败: %w", err)
	}

//获取通用数据库对象并配置连接池
  sqlDB, err := db.DB()
  if err != nil {
    return nil, err
  }
  
  sqlDB.SetMaxIdleConns(10)    // 空闲连接数
  sqlDB.SetMaxOpenConns(100)   // 最大打开连接数
  sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

  return db, nil
}