package router

import (
	"StudentService/internal/app/student"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *student.Handler) *gin.Engine {
	r := gin.Default()//创建一个服务,返回默认的路由引擎
	studentGroup := r.Group("/students")//访问地址，执行对应的方法
	{
		studentGroup.GET("", h.ListStudents)
		studentGroup.POST("", h.CreateStudent)
		studentGroup.GET("/:id", h.GetStudent)
		studentGroup.PUT("/:id", h.UpdateStudent)
		studentGroup.DELETE("/:id", h.DeleteStudent)
	}
	return r
}