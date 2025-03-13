package dao

import (
	"database/sql"
	"log"

	"github.com/zekai2023/eightball_generator/model"
)

// 查询问题的函数
func GetQuestionByID(db *sql.DB, questionID string) (*model.Question, error) {
	var question model.Question
	query := "SELECT id, question_text, correct_option FROM questions WHERE id = ?"

	// 执行查询并扫描结果
	err := db.QueryRow(query, questionID).Scan(&question.ID, &question.QuestionText, &question.CorrectOption)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有找到该问题，返回空
			return nil, nil
		}
		// 其他错误
		log.Println("查询问题失败:", err)
		return nil, err
	}

	return &question, nil
}

// 查询选项的函数
func GetOptionsByQuestionID(db *sql.DB, questionID int) ([]string, error) {
	var options []string
	query := "SELECT option_label, option_text FROM options WHERE question_id = ?"

	rows, err := db.Query(query, questionID)
	if err != nil {
		log.Println("查询选项失败:", err)
		return nil, err
	}
	defer rows.Close()

	// 遍历每行结果并将选项添加到切片中
	for rows.Next() {
		var label, text string
		if err := rows.Scan(&label, &text); err != nil {
			log.Println("扫描选项失败:", err)
			return nil, err
		}
		options = append(options, label+": "+text)
	}

	return options, nil
}
