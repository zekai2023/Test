package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	appID     = "wxd0db9418482b7149"               // 替换成你的小程序 AppID
	appSecret = "01ab10a34f556ec1df7083be91b860ac" // 替换成你的小程序 AppSecret
)

// SessionResponse 用于解析微信接口返回的数据
type SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

// LoginRequest 定义前端发送的登录请求结构
type LoginRequest struct {
	Code string `json:"code"` // 前端传来的 code
	// 可选：你也可以传递部分用户信息，如昵称、头像等
	UserInfo struct {
		NickName  string `json:"nickName"`
		AvatarUrl string `json:"avatarUrl"`
	} `json:"userInfo"`
}

// LoginResponse 定义返回给前端的结构
type LoginResponse struct {
	OpenID string `json:"openid"`
	// 如果需要返回更多用户信息，可以添加其他字段
}

// LoginHandler 处理微信登录请求
func LoginHandler(c *gin.Context) {
	// 打印前端发送的原始请求体内容
	fmt.Println("接收成功", c.Request.Body)
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read request body: %v", err)})
		fmt.Println("1111111111111111")
		return
	}

	// 打印请求体的内容，帮助调试
	fmt.Println("Received Request Body:", string(bodyBytes))

	// 重新设置请求体的读取位置，因为已经被读取过一次
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// 解析前端发送的 JSON 数据
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		fmt.Println("222222222222222222222222")
		return
	}

	// 构造调用微信接口的 URL
	wxURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, loginReq.Code,
	)

	// 调用微信接口获取 session 信息
	resp, err := http.Get(wxURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to call WeChat API: %v", err)})
		fmt.Println("333333333333333333333333")
		return
	}
	defer resp.Body.Close()

	// 读取微信返回的响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read WeChat response: %v", err)})
		fmt.Println("4444444444444444444444")
		return
	}

	// 解析微信返回的 JSON 数据
	var sessionResp SessionResponse
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse WeChat response: %v", err)})
		fmt.Println("5555555555555")
		return
	}

	// 判断微信接口是否返回错误
	if sessionResp.ErrCode != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": sessionResp.ErrMsg})
		fmt.Println("6666666666666666")
		return
	}

	// 返回 openid 给前端
	loginRes := LoginResponse{
		OpenID: sessionResp.OpenID,
	}

	// 设置响应头为 JSON 格式并返回数据
	c.JSON(http.StatusOK, loginRes)
}
