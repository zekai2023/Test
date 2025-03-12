// package main

// import (
// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	"github.com/zekai2023/eightball_generator/router"
// )

// func main() {
// 	r := gin.Default()

// 	// 配置 CORS，允许所有来源（根据需要调整）
// 	r.Use(cors.Default())

//		router.InitRoutes(r)  // 初始化路由
//		r.Run("0.0.0.0:8080") // 监听所有 IP        // 运行服务器
//	}
package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zekai2023/eightball_generator/router"
)

func main() {
	r := gin.Default()

	// 配置 CORS，允许所有来源（根据需要调整）
	r.Use(cors.Default())

	// 提供静态文件服务，确保你的 index.html 在 public 文件夹下
	r.Static("/static", "./public")

	// 初始化路由
	router.InitRoutes(r)

	// 运行服务器
	r.Run(":80") // 监听 80 端口
}
