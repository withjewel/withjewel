package utils

import "regexp"

/*VerifyUserName 验证str是否满足用户名的格式。
 */
func VerifyUsername(str string) bool {
	if str == "" {
		return false
	}
	UsernameRegex := regexp.MustCompile(`[A-Za-z0-9]{4,17}`)
	return len(UsernameRegex.FindString(str)) == len(str)
}

/*VerifyEmail 验证str是否满足电邮格式。
 */
func VerifyEmail(str string) bool {
	if str == "" {
		return false
	}
	EmailRegex := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	return len(EmailRegex.FindString(str)) == len(str)
}
