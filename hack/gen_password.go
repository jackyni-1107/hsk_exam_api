//go:build ignore
// 已迁至项目根执行: go run ./cmd/genpassword
// （单文件 go run hack/gen_password.go 在部分 GOROOT 配置下易误用标准库路径导致编译失败）

package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	h, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	fmt.Println(string(h))
}
