package common

import (
	"fmt"

	gocommon "github.com/liuhengloveyou/go-common"
	passport "github.com/liuhengloveyou/passport/client"
	"github.com/liuhengloveyou/passport/session"

	_ "github.com/go-sql-driver/mysql"
)

type AccessConfig struct {
	Addr     string                  `json:"addr"`
	Port     int                     `json:"port"`
	NodeName string                  `json:"nodeName"`
	NodeConf string                  `json:"nodeConf"`
	Passport string                  `json:"passport"`
	Session  *session.SessionManager `json:"session"`
	DBs      interface{}             `json:"dbs"`
}

type LogicConfig struct {
	NodeConf string `json:"nodeConf"`
	Nodes    []struct {
		Name  string            `json:"name"`
		Works map[string]string `json:"works"`
	} `json:"nodes"`
	Session *session.SessionManager `json:"session"`
	DBs     interface{}             `json:"dbs"`
}

var (
	AccessConf AccessConfig // 接入层配置信息
	LogicConf  LogicConfig  // 逻辑层系统配置信息

	Passport *passport.Passport
	GID      *gocommon.GlobalID
	DBs      = make(map[string]*gocommon.DBmysql)
)

func InitAccessServ(confile string) error {
	if e := gocommon.LoadJsonConfig(confile, &AccessConf); e != nil {
		return e
	}

	if nil == session.InitDefaultSessionManager(AccessConf.Session) {
		return fmt.Errorf("InitDefaultSessionManager err.")
	}

	if e := gocommon.InitDBPool(AccessConf.DBs, DBs); e != nil {
		return e
	}

	Passport = &passport.Passport{ServAddr: AccessConf.Passport}

	GID = &gocommon.GlobalID{Type: AccessConf.NodeName}

	return nil
}

func InitLogicServ(confile string) error {
	if e := gocommon.LoadJsonConfig(confile, &LogicConf); e != nil {
		return e
	}

	if nil == session.InitDefaultSessionManager(LogicConf.Session) {
		return fmt.Errorf("InitDefaultSessionManager err.")
	}

	GID = &gocommon.GlobalID{}

	return nil
}
