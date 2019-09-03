package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const HB = 3

var addr = flag.String("serv", "127.0.0.1:8080", "im server's addr.")

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println("dial ERR:", err.Error())
		return
	}
	defer conn.Close()

	// headbeat
	go func() {
		for {
			time.Sleep(time.Duration(HB) * time.Second)
			if _, err = conn.Write([]byte{0, 0, 0, 0}); err != nil {
				fmt.Println("headbeat err...")
				os.Exit(0)
			}
		}
	}()

	go func() {
		sms := make([]byte, 8192)

		c, err := conn.Read(sms)
		if err != nil {
			fmt.Println("读取服务器数据异常:", err.Error())
		}

		fmt.Println(string(sms[0:c]))

	}()

	for {
		sms := make([]byte, 8192)

		fmt.Print("请输入要发送的消息:")
		_, err := fmt.Scan(&sms)
		if err != nil {
			fmt.Println("数据输入异常:", err.Error())
			continue
		}

		_, err = conn.Write(sms)
		if err != nil {
			panic(err)
		}
	}
}
