package dao

import (
	"fmt"

	"github.com/zekai2023/eightball_generator/config"
	"github.com/zekai2023/eightball_generator/model"
)

// 根据 openid 获取错题数据
func GetWrongQuestionsByOpenID(openid string) ([]model.WrongQuestion, error) {
	// 使用原生 SQL 查询
	rows, err := config.DB.Query("SELECT * FROM wrongquestions WHERE openid = ?", openid)
	if err != nil {

		return nil, fmt.Errorf("failed to query wrong questions: %v", err)
	}
	defer rows.Close()

	// 解析查询结果
	var wrongQuestions []model.WrongQuestion
	for rows.Next() {
		var wq model.WrongQuestion
		if err := rows.Scan(&wq.ID, &wq.OpenID, &wq.Question, &wq.Answer, &wq.CreateTime, &wq.UserAnswer); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		wrongQuestions = append(wrongQuestions, wq)
	}

	// 检查是否有读取错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	return wrongQuestions, nil
}

func AddWrongQuestion(wq model.WrongQuestion) error {
	// 使用原生 SQL 进行插入
	query := "INSERT INTO wrong_questions (openid, question, answer, user_answer) VALUES (?, ?, ?, ?, ?)"

	// 执行 SQL 语句
	_, err := config.DB.Exec(query, wq.OpenID, wq.Question, wq.Answer, wq.UserAnswer)
	if err != nil {
		return fmt.Errorf("failed to insert wrong question: %v", err)
	}

	return nil
}
