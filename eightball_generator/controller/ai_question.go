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
	Category    string            `json:"category"`
}

// 全局 category 变量
var category string

// GenerateQuestion 通过 AI 生成题目
func GenerateQuestion(c *gin.Context) {
	log.Println("Received request at /generate-question")

	var requestData1 map[string]interface{} // 或者根据你的具体数据结构定义结构体
	if err := c.BindJSON(&requestData1); err != nil {
		log.Println("解析 JSON 数据失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}
	log.Println("接收到的 JSON 数据:", requestData1)
	topics := requestData1["topics"].([]interface{}) // 将 topics 转换为 []interface{}
	// 构建 AI 请求的句子
	var aiRequestContent string
	for _, topic := range topics {
		topicData := topic.(map[string]interface{})         // 将每个 topic 转换为 map[string]interface{}
		category = topicData["category"].(string)           // 获取 category 字段并设置为全局变量
		subtopics := topicData["subtopics"].([]interface{}) // 获取 subtopics 字段，它是一个数组

		// 将 subtopics 拼接成字符串
		subtopicStr := ""
		for _, subtopic := range subtopics {
			subtopicStr += subtopic.(string) + "，"
		}

		// 去掉最后的 "，"
		if len(subtopicStr) > 0 {
			subtopicStr = subtopicStr[:len(subtopicStr)-1]
		}

		// 拼接 AI 请求句子
		aiRequestContent += "根据" + category + "生成题目，包含" + subtopicStr + "；"
	}

	// 去掉最后的 "；"
	if len(aiRequestContent) > 0 {
		aiRequestContent = aiRequestContent[:len(aiRequestContent)-1]
	}

	// 输出最终的 AI 请求内容
	log.Println("构建的 AI 请求内容:", aiRequestContent)

	// 将 aiRequestContent 作为 AI 请求的内容发送

	// topic := c.DefaultQuery("topic", "A")
	log.Println("Topic received:", topics)

	// 生成 AI 请求
	aiRequest := AIRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "请根据 " + aiRequestContent + " 生成 10 道选择题，每道题包含问题、四个选项（A/B/C/D）、正确答案、解析，以及领域分类（category）。每道题的格式为一个对象，格式如下：\n{\n  \"category\": \"领域分类\",\n  \"question\": \"问题内容\",\n  \"options\": {\n    \"A\": \"选项A\",\n    \"B\": \"选项B\",\n    \"C\": \"选项C\",\n    \"D\": \"选项D\"\n  },\n  \"correct_answer\": \"正确答案\",\n  \"explanation\": \"解析内容\"\n} \n返回格式必须是JSON数组，每个问题必须包含领域分类字段,注意返回的category要是:" + category},
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
	var questionsData []QuestionData
	if err := json.Unmarshal([]byte(content), &questionsData); err != nil {
		log.Println("JSON 解析失败:", err)
		log.Println("清理后的 JSON:", content)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析题目失败"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"questions": questionsData,
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
