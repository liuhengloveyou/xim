package service

import (
	"strings"

	"github.com/liuhengloveyou/xim/common"
	"github.com/liuhengloveyou/xim/dao"
)

type User struct {
	Userid    string `validate:"-" json:"userid,omitempty"`
	Cellphone string `validate:"noneor,cellphone" json:"cellphone,omitempty"`
	Email     string `validate:"noneor,email" json:"email,omitempty"`
	Nickname  string `validate:"noneor,max=20" json:"nickname,omitempty"`
	Password  string `validate:"nonone,min=6,max=64" json:"password,omitempty"`
	Client    string `json:"client,omitempty"`
}

func (p *User) AddUser() (e error) {
	p.pretreat()

	// passport添加用户
	var userid string
	if userid, e = common.Passport.UserAdd(p.Cellphone, p.Email, p.Nickname, p.Password); e != nil {
		return e // 只能人工处理了
	}
	p.Userid = userid

	if e = p.toDao().Insert(); e != nil {
		return
	}

	return
}

func (p *User) Info() (e error) {
	return nil
}

////////
func (p *User) pretreat() {
	if p.Cellphone != "" {
		p.Cellphone = strings.ToLower(p.Cellphone)
	}
	if p.Email != "" {
		p.Email = strings.ToLower(p.Email)
	}
}

func (p *User) toDao() *dao.User {
	dao := &dao.User{Userid: p.Userid}

	return dao
}
