package service

import (
	"container/list"
	"encoding/json"
	"sync"

	"github.com/liuhengloveyou/xim/common"
)

type UserSession struct {
	ID      string
	MsgChan chan string

	inc      int
	messages *list.List
	cond     *sync.Cond
}

func NewUserSession(id string) *UserSession {
	if id == "" {
		return nil
	}

	one := &UserSession{
		ID:       id,
		MsgChan:  make(chan string, 0),
		messages: list.New()}

	one.cond = sync.NewCond(new(sync.Mutex))
	go one.popMessage()

	return one
}

func (p *UserSession) Destroy() {
	p.messages.Init()
	p.ID = ""
	close(p.MsgChan)
	p.cond.Broadcast()
}

func (p *UserSession) PushMessage(message *common.MessageForward) {
	p.messages.PushBack(message)
	p.cond.Signal()
}

func (p *UserSession) popMessage() {
	for p.ID != "" {
		p.cond.L.Lock()

		for p.ID != "" && nil == p.messages.Front() {
			p.cond.Wait()
		}

		if e := p.messages.Front(); e != nil {
			message := p.messages.Remove(e).(*common.MessageForward)
			byteMsg, _ := json.Marshal(message)
			p.MsgChan <- string(byteMsg)
		}

		p.cond.L.Unlock()
	}
}

func (p *UserSession) TopMessage() (message *common.MessageForward) {
	if e := p.messages.Front(); e != nil {
		message = e.Value.(*common.MessageForward)
	}

	return
}
