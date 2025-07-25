package main
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt" 
)


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
// // func name(name string){
// // 	fmt.Println("hello",name)
// // }



func hello(c *gin.Context){
//200是http协议响应的状态码http.StatusOk 
//gin.H是一个返回给前端的map，key是字符串，value是一个接受任意值类型的空接口 map[string]interface{}
	c.JSON(200,gin.H{"message":"Hello Go!"})
}

	type json struct{
		Name string
		Age int `json:"age"`
		Class string
	}

func main(){
	//创建一个服务,返回默认的路由引擎
	ginServer:=gin.Default()


	ginServer.GET("/Hello",hello)//用户通过浏览器使用GET请求访问“/Hello”这个地址的时候执行hello函数
	//*************gin框架返回JSON第一种使用map[string]interface{}
	//访问地址，处理我们的请求。使用RESTful API实现不同的请求执行不同的功能
	ginServer.GET("/hello",func(c *gin.Context){
		c.JSON(200,gin.H{
			"message":"hello Go!",
		})
	})
	//*************gin框架返回JSON第二种使用结构体
	//只有首字母大写的字段才会被 JSON 包导出，否则编码时直接被忽略，也可以使用struct tag（结构标签）灵活对结构体字段定制化操作,告诉 JSON 库：“把这个字段序列化成指定名字”。

	ginServer.GET("/Json",func(c *gin.Context){
		c.JSON(http.StatusOK,json{
			"小明",
			19,
			"计231",
		})
	})

	//**************获取浏览器发出请求携带的querystring参数
	ginServer.GET("querystring",func(c *gin.Context){
		name:=c.Query("query")
		age:=c.Query("age")
		c.JSON(http.StatusOK,gin.H{
			"name":name,
			"age":age,
		})
	})
	//**************获取浏览器发出请求携带的form参数
	//**************获取浏览器发出请求携带的path(URL)参数,返回字符串类型数据
	ginServer.GET("user/:name/:age",func(c *gin.Context){
		name:=c.Param("name")
		age:=c.Param("age")
		c.JSON(http.StatusOK,gin.H{
			"name":name,
			"age":age,
		})
	})
	//**************gin的JSON参数绑定
	ginServer.POST("/json", func(c *gin.Context) {
    var j json         // 声明一个 json 类型的变量 j
    err := c.ShouldBind(&j) // 把请求体里的 JSON 绑定到结构体 j
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    fmt.Printf("%#v\n", j) // 打印收到的结构体%,#v 会把任何值打印成 最详细、最精确的 Go 语法形式
    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
    })
})
	//**************gin的请求重定向
	ginServer.GET("/a",func(c *gin.Context){
		c.Request.URL.Path="/Hello"//请求的URL修改
		ginServer.HandleContext(c)
})
	ginServer.GET("/index",func(c *gin.Context){
		c.Redirect(http.StatusMovedPermanently,"http://www.baidu.com")
	})
	//**************gin的路由和路由组
	ginServer.NoRoute(func(c *gin.Context){
		c.JSON(http.StatusNotFound,gin.H{
			"msg":"此页面无效",
		})
	})
	//把共用的前缀提取出来，创建一个路由组
	videgroup:=ginServer.Group("/video")
	{
		videgroup.GET("/aaa",func(c *gin.Context){
			c.JSON(http.StatusOK,gin.H{
				"msg":"/video/aaa",
		})
	})
		videgroup.GET("/bbb",func(c *gin.Context){
			c.JSON(http.StatusOK,gin.H{
				"msg":"/video/bbb",
		})
	})
}




	//服务器端口
	ginServer.Run(":8081")
}