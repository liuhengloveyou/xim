package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/liuhengloveyou/xim/face"
	"github.com/liuhengloveyou/xim/logic"
)

var Sig string

var (
	service = flag.String("service", "", "启动什么服务? [access | logic | data]")
)

func main() {
	flag.Parse()

	sigHandler()

	switch *service {
	case "access":
		face.AccessMain()
	case "logic":
		logic.LogicMain()
	case "data":

	default:
		flag.Usage()
		os.Exit(0)
	}
}

func sigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		s := <-c
		Sig = "service is suspend ..."
		fmt.Println("Got signal:", s)
	}()
}
