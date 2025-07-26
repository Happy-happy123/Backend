package config

import (
	"strconv"
	"os"
)

// type Config struct {
// 	Port        int
// 	DatabaseURL string
// }

// func Load() Config {
// 	port, err := strconv.Atoi(getenv("PORT","8080")) //将有效整数的字符串转换为整数
// 	if err!= nil{
// 		fmt.Println("转换失败:",err)
// 		return Config{
// 			Port:8080,
// 			DatabaseURL:getenv("DATABASE_URL", "gouser:Mxd20051020@@tcp(127.0.0.1:3306)/student_service"),
// 		}
// 	}else{
// 		return Config{
// 			Port:port,
// 			DatabaseURL:getenv("DATABASE_URL","gouser:Mxd20051020@@tcp(127.0.0.1:3306)/student_service"),
// 		}
// 	}
// }
type Config struct {
  Port        int
  MySQLDSN    string
  RedisAddr   string
  RedisDB     int
}

//获取环境变量或默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// 获取整数环境变量
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func Load() Config {
	return Config{
		Port:       getEnvAsInt("PORT", 8080),
		MySQLDSN: getEnv("MYSQL_DSN", "gouser:Mxd20051020@@tcp(127.0.0.1:3306)/student_service"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:     getEnvAsInt("REDIS_DB", 0),
	}
}



