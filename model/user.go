package model

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID                 int
	Username, Password string
	Email              sql.NullString
}

func Verify(username string, password string) bool {
	db, err := gorm.Open("mysql", "leslie:19941121sr@tcp(115.28.27.5:3306)/test")
	if err != nil {
		fmt.Printf("数据库连接出错：%s\n", err.Error())
		return false
	}
	db.SingularTable(true)
	defer db.Close()

	var user User
	return !db.Find(&user, "username = ? and password = ?", username, password).RecordNotFound()
}

func VerifyEmail(email string, password string) (bool, bool) {
	db, err := gorm.Open("mysql", "leslie:19941121sr@tcp(115.28.27.5:3306)/test")
	if err != nil {
		fmt.Printf("数据库连接出错：%s\n", err.Error())
		return false, false
	}
	db.SingularTable(true)
	defer db.Close()

	var user User
	exist := !db.Find(&user, "email = ?", email).RecordNotFound()
	if !exist {
		return false, false
	}
	return true, user.Password == password
}
