package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// 声明一个全局变量 db，用于数据库连接
var DB *sql.DB

// 初始化数据库连接
func InitDB() {
	var err error
	// 配置数据库连接字符串：用户名:密码@协议(主机:端口)/数据库名
	DB, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/eightball_generator")
	if err != nil {
		log.Fatal(err)
	}
	// 测试数据库连接
	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("数据库连接成功")
}
