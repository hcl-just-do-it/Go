package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func main() {
	//http连接信息打印到gin_http.log文件里而不是控制台，每次重启服务删除了之前日志
	logfile, err := os.Create("./gin_http.log")
	if err != nil {
		fmt.Println("Could not create log file")
	}
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = io.MultiWriter(logfile)
	//
	r := gin.Default()

	initRouter(r) //first

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
