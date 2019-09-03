package logic

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
)

func init() {
	nodenet.RegisterWorker("TempGroupLogin", common.MessageTGLogin{}, TempGroupLogin)
	nodenet.RegisterWorker("TempGroupSend", common.MessageForward{}, TempGroupSend)
}

func TempGroupLogin(data interface{}) (result interface{}, err error) {
	var msg = data.(common.MessageTGLogin)
	log.Infoln("tgroupLogin:", msg)

	sess, err := session.GetSessionById(msg.Gid)
	if err != nil {
		return nil, err
	}

	if err := sess.Set(msg.Uid, msg.Access); err != nil {
		return nil, err
	}

	log.Infoln("tgroupLogin OK:", sess)

	return nil, nil
}

func TempGroupSend(data interface{}) (result interface{}, err error) {
	var msg = data.(common.MessageForward)
	log.Infoln(msg)

	sess, err := session.GetSessionById(msg.ToUserid)
	if err != nil {
		return nil, err
	}

	keys := sess.Keys()
	log.Infoln("tempGroupSend:", msg.ToUserid, keys)
	for i := 0; i < len(keys); i++ {
		stat := sess.Get(keys[i])
		if stat == nil || msg.FromUserid == keys[i] {
			log.Errorln("tgroup skip:", keys[i], stat)
			continue
		}

		cMsg := nodenet.NewMessage("", "", nil, common.MessageForward{FromUserid: msg.FromUserid, ToUserid: fmt.Sprintf("%v.%v", msg.ToUserid, keys[i]), ToGroupId: msg.ToUserid, Content: msg.Content})
		log.Infoln("tgroup pushmsg: ", stat.(string), cMsg)

		if err = nodenet.SendMsgToComponent(stat.(string), cMsg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
