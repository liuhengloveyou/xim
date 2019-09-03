/*
 * 消息路由
 */

package logic

import (
	"container/list"
	"time"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
)

func ForwardMessage(data interface{}) (result interface{}, err error) {
	var msg = data.(common.MessageForward)
	log.Infof("ForwardMessage: %#v\n", msg)
	newid := common.GID.LogicClock(msg.MsgId)
	log.Infof("logicClock: %d->%d\n", msg.MsgId, newid)
	msg.MsgId = newid // 全局排序

	sess, err := session.GetSessionById(msg.ToUserid)
	if err != nil {
		return nil, err
	}
	log.Infof("ForwardMessage tosession: %v, %#v\n", msg.ToUserid, sess)

	if sess.Get("info") == nil {
		sess.Set("info", &StateSession{})
	}
	info := sess.Get("info").(*StateSession)
	setOfflineMessage(info, &msg) // 先放到离线消息队列

	msg.FromeAccess = ""
	msg.Time = time.Now().Unix()

	if info.Alive <= 0 {
		log.Infof("offline: %v. \n", msg.ToUserid)
		return nil, nil
	}

	dealOfflineMessage(sess, info) // 处理离线消息

	return nil, nil
}

func dealOfflineMessage(sess session.SessionStore, info *StateSession) {
	info.lock.Lock()
	defer info.lock.Unlock()

	// 删除已经确认的
	for info.Messages.Len() > 0 {
		el := info.Messages.Front()
		one := el.Value.(*common.MessageForward)
		if one.MsgId >= info.Confirm {
			break
		}

		info.Messages.Remove(el)
	}

	for e, max := info.Messages.Front(), 5; e != nil && max > 0; e, max = e.Next(), max-1 {
		one := e.Value.(*common.MessageForward)
		log.Infoln("dealOfflineMessage:", one)
		if one.MsgId <= info.Pushed {
			continue
		}

		log.Infof("pushmsg: %#v. %#v", info, one)
		if pushMessage(sess, one) == true {
			info.Pushed = one.MsgId // 推送成功
		}
	}

}

// 一次推一条
func pushMessage(sess session.SessionStore, message *common.MessageForward) (ok bool) {
	cMsg := nodenet.NewMessage("", "", make([]string, 0), nil)
	keys := sess.Keys()

	for i := 0; i < len(keys); i++ {
		if keys[i] == "info" {
			continue
		}
		if sess.Get(keys[i]) == nil {
			continue
		}
		stat := sess.Get(keys[i]).(*common.MessageLogin)
		if stat.UpdateTime <= 0 {
			continue // 已退出
		}

		message.ToAccess = stat.AccessName
		message.ToSession = stat.AccessSession
		cMsg.Payload = message
		log.Infof("ForwardMessage: %s, %#v, %#v", keys[i], stat, message)

		if err := nodenet.SendMsgToComponent(stat.AccessName, cMsg); err != nil {
			log.Infof("ForwardMessage ERR: %s, %#v, %#v; %v", keys[i], stat, message, err)
		} else {
			ok = true
		}
	}

	return
}

func setOfflineMessage(info *StateSession, message *common.MessageForward) {
	info.lock.Lock()
	defer info.lock.Unlock()

	if info.Messages == nil {
		info.Messages = list.New()
	}

	log.Infof("setOfflineMessage: %#v. %#v", info, message)
	info.Messages.PushBack(message)
}
