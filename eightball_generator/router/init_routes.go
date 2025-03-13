package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zekai2023/eightball_generator/controller"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	r.GET("/ping", controller.Ping) // 测试接口
	r.POST("/ask", controller.Ask)  // 新增 /ask 接口
	r.GET("/question", controller.Question_Get)
	r.GET("/question1", controller.GetQuestion)
}
