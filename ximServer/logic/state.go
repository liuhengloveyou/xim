/*
* 用户长连接在线状态更新
 */
package logic

import (
	"container/list"
	"sync"
	"time"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
	"github.com/liuhengloveyou/passport/session"
)

type StateSession struct {
	Alive    int64      // 最后活动时间
	Pushed   int64      // 已经推送消息指针
	Confirm  int64      // 已经确认消息指针
	Messages *list.List // 离线消息
	lock     sync.Mutex
}

func init() {
	nodenet.RegisterWorker("UerLogin", common.MessageLogin{}, UerLogin)
	nodenet.RegisterWorker("UerLogout", common.MessageLogout{}, UerLogout)
	nodenet.RegisterWorker("ForwardMessage", common.MessageForward{}, ForwardMessage)
}

func UerLogin(data interface{}) (result interface{}, err error) {
	var msg = data.(common.MessageLogin)
	log.Infoln(msg)
	if msg.ClientType == "" {
		msg.ClientType = "XIM"
	}

	sess, err := session.GetSessionById(msg.Userid)
	if err != nil {
		return nil, err
	}

	// 互踩@@

	// 状态会话
	if nil == sess.Get("info") {
		sess.Set("info", &StateSession{})
	}
	info := sess.Get("info").(*StateSession)

	msg.UpdateTime = time.Now().Unix()
	if err = sess.Set(msg.ClientType, &msg); err != nil {
		return nil, err
	}
	log.Infof("UserLogin OK: %#v; %#v", msg, sess)

	// 最后活动时间
	info.Alive = time.Now().Unix()

	// 推送消息确认
	if msg.ConfirmMessage > info.Confirm {
		info.Confirm = msg.ConfirmMessage
	}

	// 离线消息
	if nil != info.Messages {
		dealOfflineMessage(sess, info) // 处理离线消息
	}

	return nil, nil
}

func UerLogout(data interface{}) (result interface{}, err error) {
	var msg = data.(common.MessageLogout)
	log.Infof("%#v", msg)

	sess, err := session.GetSessionById(msg.Userid)
	if err != nil {
		return nil, err
	}
	log.Infof("current session: %#v; %#v", msg, sess)

	if nil == sess.Get(msg.ClientType) {
		log.Errorf("UerLogout sess nil: %#v.", msg)
		return nil, nil
	}

	stat := sess.Get(msg.ClientType).(*common.MessageLogin)
	if stat.AccessName != msg.AccessName || stat.AccessSession != msg.AccessSession {
		log.Errorf("UerLogout sess ERR: %#v | %#v", msg, stat)
		return nil, nil
	}

	stat.UpdateTime = -1 // 不在线状态

	return nil, nil
}

// 踩人
func kickOff() {

}
