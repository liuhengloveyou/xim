package service

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
)

func TGroutRecv(uid, gid string) (info *UserSession, e error) {
	if uid == "" {
		return nil, fmt.Errorf("userid nil")
	}
	if gid == "" {
		return nil, fmt.Errorf("groupid nil")
	}

	sid := fmt.Sprintf("%s.%s", gid, uid)
	sess, e := session.GetSessionById(sid)
	if e != nil {
		return nil, e
	}

	if sess.Get("info") != nil {
		log.Infoln("tgroup userlogined:", gid, uid)
		return sess.Get("info").(*UserSession), nil
	}

	log.Infoln("tgroup userlogin:", gid, uid)
	sess.Set("info", NewUserSession(fmt.Sprintf("%s.%s", gid, uid)))

	g := nodenet.GetGraphByName(common.LOGIC_TEMPGROUP)
	if len(g) < 1 {
		return nil, fmt.Errorf("graph nil:", common.LOGIC_TEMPGROUP)
	}

	cMsg := nodenet.NewMessage(common.GID.ID(), common.AccessConf.NodeName, g, common.MessageTGLogin{Uid: uid, Gid: gid, Access: common.AccessConf.NodeName})
	cMsg.DispenseKey = gid

	if e = nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return nil, e
	}

	return info, nil
}

func TGroutSend(uid, gid, message string) error {
	cMsg := nodenet.NewMessage(common.GID.ID(), common.AccessConf.NodeName, nodenet.GetGraphByName(common.LOGIC_TEMPGROUP), &common.MessageForward{FromUserid: uid, ToUserid: gid, Content: message})
	cMsg.DispenseKey = gid
	log.Infoln(cMsg)

	if e := nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return e
	}

	return nil
}
