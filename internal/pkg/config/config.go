package config

import (
	"fmt"
	"strconv"
	"os"
)

type Config struct {
	Port        int
	DatabaseURL string
}

func Load() Config {
	port, err := strconv.Atoi(getenv("PORT","8080")) //将有效整数的字符串转换为整数
	if err!= nil{
		fmt.Println("转换失败:",err)
		return Config{
			Port:8080,
			DatabaseURL:getenv("DATABASE_URL", "gouser:Mxd20051020@@tcp(127.0.0.1:3306)/student_service"),
		}
	}else{
		return Config{
			Port:port,
			DatabaseURL:getenv("DATABASE_URL","gouser:Mxd20051020@@tcp(127.0.0.1:3306)/student_service"),
		}
	}
}

func getenv(key,defaultvalue string) string{
	value,exists:=os.LookupEnv(key)//从环境变量中查找指定的键（key）并获取其对应的值
	if exists{
		return value
	}else{
		return defaultvalue
	}

}

