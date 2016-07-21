package model

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	id                 int
	username, password string
	email              sql.NullString
}

func Verify(username string, password string) bool {
	db, err := sql.Open("mysql", "leslie:19941121sr@tcp(115.28.27.5:3306)/test")
	if err != nil {
		fmt.Printf("数据库连接出错：%s\n", err.Error())
		return false
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("ping不通数据库：%s\n", err.Error())
		return false
	}

	result := db.QueryRow("select * from user where username = ? and password = ?", username, password)
	if result == nil {
		fmt.Printf("在查询的时候出错\n")
		return false
	}
	var user User
	err = result.Scan(&user.id, &user.username, &user.password, &user.email)
	if err != nil {
		return false
	}

	fmt.Println(user)
	return true
}

func VerifyEmail(email string, password string) bool {
	db, err := sql.Open("mysql", "leslie:19941121sr@tcp(115.28.27.5:3306)/test")
	if err != nil {
		fmt.Printf("数据库连接出错：%s\n", err.Error())
		return false
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("ping不通数据库：%s\n", err.Error())
		return false
	}

	result := db.QueryRow("select * from user where email = ? and password = ?", email, password)
	if result == nil {
		fmt.Printf("在查询的时候出错\n")
		return false
	}
	var user User
	err = result.Scan(&user.id, &user.username, &user.password, &user.email)
	if err != nil {
		return false
	}

	fmt.Println(user)
	return true
}
