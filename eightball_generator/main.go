package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zekai2023/eightball_generator/config"
	"github.com/zekai2023/eightball_generator/router"
)

func main() {
	config.InitDB()
	r := gin.Default()

	// 配置 CORS，允许所有来源（根据需要调整）
	r.Use(cors.Default())

	router.InitRoutes(r)  // 初始化路由
	r.Run("0.0.0.0:8080") // 监听所有 IP        // 运行服务器
}
