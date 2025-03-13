package model

// Question 结构体用于存储问题和选项
type Question struct {
	ID            int    `json:"id"`
	QuestionText  string `json:"question_text"`
	CorrectOption string `json:"correct_option"` // 添加这个字段来存储正确答案
}
