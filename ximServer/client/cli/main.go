package main

import (
	"fmt"
)

func main() {
	sessionID := ""
	phone := ""
	fmt.Println("临时会话编号:")
	if _, e := fmt.Scan(&sessionID); e != nil {
		panic(e)
	}
	fmt.Println("您的手机号码:")
	if _, e := fmt.Scan(&phone); e != nil {
		panic(e)
	}
	fmt.Println(phone, "加入了临时会话", sessionID)
}
