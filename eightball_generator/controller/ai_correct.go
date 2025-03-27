package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 用户提交的答案结构
type UserAnswer struct {
	Question   string `json:"question"`
	UserAnswer string `json:"userAnswer"`
}

// 批改结果结构
type GradingResult struct {
	Question    string `json:"question"`
	UserAnswer  string `json:"userAnswer"`
	IsCorrect   bool   `json:"isCorrect"`
	Explanation string `json:"explanation"`
}

type SubmissionRequest struct {
	Submissions []UserAnswer `json:"submissions"`
}

// 批改答案的接口
func GradeAnswers(c *gin.Context) {
	var submissionRequest SubmissionRequest
	if err := c.ShouldBindJSON(&submissionRequest); err != nil {
		log.Println("解析请求数据失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	var gradingResults []GradingResult

	// 批改每道题
	for _, userAnswer := range submissionRequest.Submissions {
		// 调用AI批改
		gradingResult := gradeAnswerWithAI(userAnswer)
		gradingResults = append(gradingResults, gradingResult)
	}

	// 返回批改结果
	c.JSON(http.StatusOK, gin.H{
		"gradingResults": gradingResults,
	})
}

// 使用AI批改答案
func gradeAnswerWithAI(userAnswer UserAnswer) GradingResult {
	// 生成AI请求
	aiRequest := AIRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "请根据以下问题和用户的答案批改并返回结果。\n问题： " + userAnswer.Question + "\n用户答案： " + userAnswer.UserAnswer + "\n请判断答案是否正确，并提供解析，返回格式如下：\n{\n  \"isCorrect\": true/false,\n  \"explanation\": \"解析内容\"\n}"},
		},
	}

	// 编码请求数据
	requestData, err := json.Marshal(aiRequest)
	if err != nil {
		log.Println("编码请求数据失败:", err)
		return GradingResult{}
	}

	// 发送请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestData))
	if err != nil {
		log.Println("创建请求失败:", err)
		return GradingResult{}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("调用 AI API 失败:", err)
		return GradingResult{}
	}
	defer resp.Body.Close()

	// 解析响应体
	var aiResponse AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResponse); err != nil {
		log.Println("解析 AI 响应失败:", err)
		return GradingResult{}
	}

	if len(aiResponse.Choices) == 0 {
		log.Println("AI 生成的数据为空")
		return GradingResult{}
	}

	// 提取AI返回的批改结果
	content := aiResponse.Choices[0].Message.Content
	log.Println("AI 批改响应内容:", content)

	// 解析AI返回的批改结果
	var gradingResult GradingResult
	gradingResult.Question = userAnswer.Question
	gradingResult.UserAnswer = userAnswer.UserAnswer
	gradingResult.Explanation = content

	// 假设判断答案是否正确基于解析中的逻辑
	if strings.Contains(content, "正确") {
		gradingResult.IsCorrect = true
	} else {
		gradingResult.IsCorrect = false
	}

	// 返回批改结果，包含问题、用户答案、批改内容
	return gradingResult
}
