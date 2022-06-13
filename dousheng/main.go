package main

import (
	"dousheng/config"
	"dousheng/database"
	"dousheng/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	cfg, err := config.ParseConfig("./config/init.json")
	if err != nil {
		panic("Configuration file parse error! " + err.Error())
	}
	// 连接Mysql
	if err := database.InitMySQL(cfg); err != nil {
		panic("MySQL connect error!" + err.Error())
	}
	sqlDB, _ := database.MySQLDB.DB()
	defer sqlDB.Close()
	// 连接redis
	if err := database.InitRedisClient(cfg); err != nil {
		panic("Reids connect error!" + err.Error())
	}
	// 连接MQ

	// 初始化路由
	r := gin.Default()
	router.InitRouter(r)

	r.Run(cfg.AppPort)
}
