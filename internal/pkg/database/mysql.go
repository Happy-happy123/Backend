package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)
var db *sql.DB
var err error
func NewMySQLClient(s string)(*sql.DB,error){
	if db,err=sql.Open("mysql",s);err!=nil{// Open不会检验用户名和密码
		return nil,fmt.Errorf("sql.open failed:%v",err)
	}
	if err = db.Ping(); err != nil {// Ping会检验用户名和密码
		return nil,fmt.Errorf("db.ping failed:%v",err)
	}
	if _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS students (
            id     INT AUTO_INCREMENT PRIMARY KEY   COMMENT '编号',
            name   VARCHAR(50)            			COMMENT '姓名',
            age    TINYINT UNSIGNED       			COMMENT '年龄',
            gender CHAR(1)               	 		COMMENT '性别'
        ) ENGINE=InnoDB
          DEFAULT CHARSET=utf8mb4
          COLLATE=utf8mb4_0900_ai_ci
          COMMENT='学生信息表';
    `);err!=nil{
		return nil,fmt.Errorf("学生表创建失败:%v",err)
	}
	return db,nil
}