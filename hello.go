// package main
// import "fmt"
// func main() {
// 	//不带声明格式（初始化声明）只能在函数体中出现
// 	// a:=19
// 	// b:="冒险岛"
// 	// c:=fmt.Sprintf("name=%s,age=%d",b,a)
// 	// fmt.Println(c)
// 	// fmt.Printf("name=%s,age=%d\n", b, a)
// 	// e,f:=1,2
// 	// fmt.Println(&e,&f)
// 	//局部变量声明后必须在局部使用，全局变量是允许声明但不使用的
// 	// g:=3
// 	// _,numb,strs := numbers() //只获取函数返回值的后两个
//   	// fmt.Println(numb,strs)


// 	// my_const()

// //正常函数调用
// 	name("go")
// //闭包一
// 	func (name string){
// 		fmt.Println("hello",name)
// }("go")
// //闭包二
// 	res:=func(name string){
// 		fmt.Println("hello",name)
// 	}
// 	res("go")
// }
// // func numbers()(int,int,string){
// //   	a , b , c := 1 , 2 , "str"
// //   	return a,b,c
// // }

// //常量
// // func my_const(){
// // 	    const (
// //             a1 = iota   //0
// //             b1          //1
// //             c1          //2
// //             d1 = "ha"   //独立值，iota += 1
// //             e1          //"ha"   iota += 1
// //             f1 = 100    //iota +=1
// //             g1          //100  iota +=1
// //             h1 = iota   //7,恢复计数
// //             i1          //8
// //     )
// //     fmt.Println(a1,b1,c1,d1,e1,f1,g1,h1,i1)
// // }

// //匿名函数
// func name(name string){
// 	fmt.Println("hello",name)
// }
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/dbname")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("连接失败:", err)
    }
    fmt.Println("已连接 MySQL！")
}