// package router

// import (
// 	"StudentService/internal/app/student"
// 	"github.com/gin-gonic/gin"
// )

// func SetupRouter(h *student.Handler) *gin.Engine {
// 	r := gin.Default()//创建一个服务,返回默认的路由引擎
// 	studentGroup := r.Group("/students")//访问地址，执行对应的方法
// 	{
// 		studentGroup.GET("", h.ListStudents)
// 		studentGroup.POST("", h.CreateStudent)
// 		studentGroup.GET("/:id", h.GetStudent)
// 		studentGroup.PUT("/:id", h.UpdateStudent)
// 		studentGroup.DELETE("/:id", h.DeleteStudent)
// 	}
// 	return r
// }

package router

import (
	"StudentService/internal/app/student"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// 访问日志中间件
func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		
		// 处理请求
		c.Next()
		
		// 记录日志
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		
		logrus.WithFields(logrus.Fields{
			"path":    path,
			"query":   query,
			"method":  method,
			"ip":      clientIP,
			"latency": latency,
			"status":  statusCode,
		}).Info("API访问日志")
	}
}

func SetupRouter(h *student.Handler) *gin.Engine {
	r := gin.Default()
	
	// 添加中间件
	r.Use(CorsMiddleware())
	r.Use(AccessLogMiddleware())
	
	studentGroup := r.Group("/students")
	{
		studentGroup.GET("", h.ListStudents)
		studentGroup.POST("", h.CreateStudent)
		studentGroup.GET("/:id", h.GetStudent)
		studentGroup.PUT("/:id", h.UpdateStudent)
		studentGroup.DELETE("/:id", h.DeleteStudent)
	}
	return r
}