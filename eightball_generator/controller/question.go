package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping 测试接口
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// Ask 处理用户输入的问题
func Ask(c *gin.Context) {
	var request struct {
		Question string `json:"question"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 返回一个固定的回答，后续可以接 AI 处理
	c.JSON(http.StatusOK, gin.H{"answer": "这是一个自动生成的回答：" + request.Question})
}
