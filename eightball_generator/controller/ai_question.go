package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// AI API 配置
const (
	apiKey = "sk-ccae2beca2d343a295e59e078c470a4f"
	apiURL = "https://api.deepseek.com/v1/chat/completions"
)

// AI请求结构体
type AIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AI返回的数据结构
type AIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// 题目结构
type QuestionData struct {
	Question    string            `json:"question"`
	Options     map[string]string `json:"options"` // **修改这里，支持对象格式**
	Answer      string            `json:"correct_answer"`
	Explanation string            `json:"explanation"`
}

// GenerateQuestion 通过 AI 生成题目
func GenerateQuestion(c *gin.Context) {
	log.Println("Received request at /generate-question")

	topic := c.DefaultQuery("topic", "A")
	log.Println("Topic received:", topic)

	// 生成 AI 请求
	aiRequest := AIRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "请根据 " + topic + " 生成一个选择题，包括问题、四个选项（A/B/C/D），正确答案和解析，使用JSON格式返回。"},
		},
	}

	// 编码请求数据
	requestData, err := json.Marshal(aiRequest)
	if err != nil {
		log.Println("编码请求数据失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成题目失败"})
		return
	}

	// 发送请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestData))
	if err != nil {
		log.Println("创建请求失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成题目失败"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("调用 AI API 失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成题目失败"})
		return
	}
	defer resp.Body.Close()

	// 解析响应体
	var aiResponse AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResponse); err != nil {
		log.Println("解析 AI 响应失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析题目失败"})
		return
	}

	// 提取 AI 返回的文本
	if len(aiResponse.Choices) == 0 {
		log.Println("AI 生成的数据为空")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 生成的数据为空"})
		return
	}
	content := aiResponse.Choices[0].Message.Content
	log.Println("AI 原始响应内容:", content)

	// **去掉 Markdown 代码块**
	content = cleanJSONContent(content)

	// 解析 JSON
	var questionData QuestionData
	if err := json.Unmarshal([]byte(content), &questionData); err != nil {
		log.Println("JSON 解析失败:", err)
		log.Println("清理后的 JSON:", content)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析题目失败"})
		return
	}

	// **转换 options 格式**
	optionsList := []string{}
	for key, value := range questionData.Options {
		optionsList = append(optionsList, key+". "+value)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"question":    questionData.Question,
		"options":     optionsList,
		"answer":      questionData.Answer,
		"explanation": questionData.Explanation,
	})
}

// cleanJSONContent 处理 JSON 字符串，去掉 ```json ... ``` 代码块
func cleanJSONContent(text string) string {
	// 去掉 Markdown 代码块
	re := regexp.MustCompile("(?s)```json\\s*(.*?)\\s*```")
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return text
}
