package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zekai2023/eightball_generator/dao"
)

// 获取错题本数据接口
func GetWrongQuestionsHandler(c *gin.Context) {
	// 获取请求参数中的 openid
	openid := c.DefaultQuery("openid", "")
	if openid == "" {
		// 如果没有提供 openid，返回错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "openid is required"})
		return
	}

	// 调用 DAO 层方法获取错题数据
	wrongQuestions, err := dao.GetWrongQuestionsByOpenID(openid)
	if err != nil {
		// 如果查询出错，返回错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果查询成功，返回数据
	c.JSON(http.StatusOK, wrongQuestions)
}
