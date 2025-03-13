package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zekai2023/eightball_generator/config"
	"github.com/zekai2023/eightball_generator/dao"
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

func Question_Get(c *gin.Context) {
	// 假设这是从数据库获取的题目数据
	question := "Go语言的创始人是谁？"
	options := []string{"A. Rob Pike", "B. Guido van Rossum", "C. Dennis Ritchie", "D. Ken Thompson"}

	c.JSON(http.StatusOK, gin.H{
		"question": question,
		"options":  options,
	})
}

func GetQuestion(c *gin.Context) {
	// 假设从请求参数中获取问题的 ID
	questionID := c.DefaultQuery("question_id", "1") // 默认取问题 ID 为 1，如果没有传入

	// 查询问题
	question, err := dao.GetQuestionByID(config.DB, questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询问题失败"})
		return
	}
	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "问题不存在"})
		return
	}

	// 查询选项
	options, err := dao.GetOptionsByQuestionID(config.DB, question.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询选项失败"})
		return
	}

	// 返回问题和选项
	c.JSON(http.StatusOK, gin.H{
		"question": question.QuestionText,
		"options":  options,
	})
}
