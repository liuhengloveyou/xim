/*
 * 消息路由
 */
package service

import (
	"encoding/json"
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
)

func SendMessage(sess session.SessionStore, body []byte) (e error) {
	var user *User
	if nil != sess.Get("user") {
		user = sess.Get("user").(*User)
	}
	if nil == user {
		return fmt.Errorf("会话错误.")
	}

	message := &common.MessageForward{}
	if e = json.Unmarshal(body, message); e != nil {
		return e
	}
	log.Infoln("SendMessage:", string(body), message)

	if message.ToUserid == "" || message.ShowType == "" || message.Content == "" {
		return fmt.Errorf("消息格式错误.")
	}
	message.FromUserid = user.Userid
	message.FromeAccess = common.AccessConf.NodeName
	message.MsgId = common.GID.LogicClock(0)

	cMsg := nodenet.NewMessage(fmt.Sprintf("%v", message.MsgId),
		common.AccessConf.NodeName,
		nodenet.GetGraphByName(common.LOGIC_FORWARD),
		message)
	cMsg.DispenseKey = message.ToUserid
	log.Infoln(cMsg, cMsg.Payload)

	if e := nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return e
	}

	return nil
}
