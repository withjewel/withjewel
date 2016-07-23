package model

import (
	"withjewel/jewel/jedb"
	"database/sql"
)


type User struct {
	ID                 int
	Username, Password string
	Email              sql.NullString
}

/*Verify 验证凭此用户名和密码是否能够登陆。
 */
func Verify(username string, password string) bool {
	db := jedb.Accquire()
	defer jedb.Revert(db)

	var user User
	return !db.Find(&user, "username = ? and password = ?", username, password).RecordNotFound()
}

/*VerifyEmail 验证凭此邮箱和密码是否能够登陆。
第一个返回值表示邮箱是否存在。
若存在，第二个返回值表示密码是否匹配。若不存在，第二个返回值表示该邮箱是否正在等待确认。如下表：
---------------------------------------------------------------
first-bool  | true           | false
---------------------------------------------------------------
second-bool
---------------------------------------------------------------
true        | 可以登录        | 邮箱已经申请注册，但是还未通过确认
---------------------------------------------------------------
false       | 邮箱与密码不匹配 | 邮箱不存在账号系统中
---------------------------------------------------------------
*/
func VerifyEmail(email string, password string) (bool, bool) {
	db := jedb.Accquire()
	defer jedb.Revert(db)

	var user User
	exist := !db.Find(&user, "email = ?", email).RecordNotFound()
	if !exist {
		return false, false
	}
	return true, user.Password == password
}
