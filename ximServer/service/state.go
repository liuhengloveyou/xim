/*
 * 用户长连接在线状态
 */
package service

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
)

func StateUpdate(sess session.SessionStore) (info *UserSession, e error) {
	if sess.Get("user") == nil {
		log.Errorln("session ERR: ", sess)
		return nil, fmt.Errorf("会话错误.")
	}
	user := sess.Get("user").(*User)

	if sess.Get("info") != nil {
		info = sess.Get("info").(*UserSession)
		info.inc += 1
		if info.inc < 3 {
			log.Infoln("state need not update:", sess)
			return sess.Get("info").(*UserSession), nil
		}
	} else {
		log.Infoln("StateUpdate:", sess)
		info = NewUserSession(sess.Id(""))
		sess.Set("info", info)
	}

	g := nodenet.GetGraphByName(common.LOGIC_STATE)
	if len(g) < 1 {
		return nil, fmt.Errorf("graph not config:", common.LOGIC_STATE)
	}

	msg := &common.MessageLogin{
		Userid:        user.Userid,
		ClientType:    user.Client,
		AccessName:    common.AccessConf.NodeName,
		AccessSession: sess.Id("")}

	// 推送消息确认
	top := info.TopMessage()
	if top != nil {
		msg.ConfirmMessage = top.MsgId
		log.Infoln("StateUpdate msg:", msg)
	}

	cMsg := nodenet.NewMessage(common.GID.ID(), common.AccessConf.NodeName, g, msg)
	cMsg.DispenseKey = user.Userid

	log.Infoln("StateUpdate send:", cMsg, cMsg.Payload)
	if e = nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return nil, e
	}

	return info, nil
}
