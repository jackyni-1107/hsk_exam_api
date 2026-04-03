package utility

import "golang.org/x/crypto/bcrypt"

// CheckPassword 校验明文是否与 bcrypt 哈希匹配。
func CheckPassword(passwordHash, plainPassword string) bool {
	if passwordHash == "" || plainPassword == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword)) == nil
}
