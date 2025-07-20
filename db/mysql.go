package db

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    // "StudentService/models"
)

var DB *sql.DB

func InitDB() {
    dsn := "gouser:Mxd20051020@@tcp(127.0.0.1:3306)/student_service"
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        fmt.Printf("dsn:%s invalid,err:%v\n", dsn, err)
        return
    }
    err = DB.Ping()
    if err != nil {
        fmt.Printf("open %s faild,err:%v\n", dsn, err)
        return
    }
    fmt.Println("连接数据库成功")
}

func CreateStudentTable() {
    _, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS student (
            id     INT PRIMARY KEY         COMMENT '编号',
            name   VARCHAR(50)            COMMENT '姓名',
            age    TINYINT UNSIGNED       COMMENT '年龄',
            gender CHAR(1)                COMMENT '性别'
        ) ENGINE=InnoDB
          DEFAULT CHARSET=utf8mb4
          COLLATE=utf8mb4_0900_ai_ci
          COMMENT='学生信息表';
    `)
    if err != nil {
        fmt.Printf("create table error: %v\n", err)
    }
    fmt.Println("学生信息表创建完成")
}

