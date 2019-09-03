/*
业务逻辑
*/

package logic

import (
	"flag"
	"fmt"
	"reflect"
	"time"

	"github.com/liuhengloveyou/xim/common"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/client"
)

var (
	Sig string

	mynodes  map[string]*nodenet.Component
	passport *client.Passport
)

var (
	confile = flag.String("c", "example/logic.conf.simple", "配置文件路径.")
)

func init() {
	mynodes = make(map[string]*nodenet.Component)
}

func initNodenet(fn string) error {
	if e := nodenet.BuildFromConfig(fn); e != nil {
		return e
	}

	for i := 0; i < len(common.LogicConf.Nodes); i++ {
		name := common.LogicConf.Nodes[i].Name
		mynodes[name] = nodenet.GetComponentByName(name)
		if mynodes[name] == nil {
			return fmt.Errorf("No node: %v.", name)
		}

		for k, v := range common.LogicConf.Nodes[i].Works {
			t, w := nodenet.GetMessageTypeByName(k), nodenet.GetWorkerByName(v)
			if t == nil {
				return fmt.Errorf("No message registerd: %s", k)
			}
			if w == nil {
				return fmt.Errorf("No worker registerd: %s", v)
			}
			if reflect.TypeOf(t) != w.Message {
				return fmt.Errorf("worker can't recive message: %v %v", w, k)
			}
			mynodes[name].RegisterHandler(t, w.Handler)
		}

		go mynodes[name].Run()
	}

	return nil
}

func LogicMain() {
	if e := common.InitLogicServ(*confile); e != nil {
		panic(e)
	}

	if e := initNodenet(common.LogicConf.NodeConf); e != nil {
		panic(e)
	}

	fmt.Println("logic GO...")
	for {
		time.Sleep(3 * time.Second)
	}
}
