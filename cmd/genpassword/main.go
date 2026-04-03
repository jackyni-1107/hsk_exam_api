// 生成 bcrypt 密码哈希，用于写入 sys_user / client_user / sys_member 等表。
// 用法（在项目根目录）: go run ./cmd/genpassword
// 修改下方 plain 常量后运行。
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	const plain = "123456"
	cost := bcrypt.DefaultCost
	if len(os.Args) > 1 {
		// 可选: go run ./cmd/genpassword mypass
		h, err := bcrypt.GenerateFromPassword([]byte(os.Args[1]), cost)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(h))
		return
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(h))
}
