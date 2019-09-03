package dao

import (
	"time"

	"github.com/liuhengloveyou/xim/common"
)

type User struct {
	Userid  string
	Icon    string
	AddTime time.Time
	Version int
}

func (p *User) Insert() (e error) {
	_, e = common.DBs["xim"].Insert("INSERT INTO user(userid,version) values(?,?)", p.Userid, time.Now().Unix())

	return
}

func (p *User) Update() (e error) {
	// _, e = common.DBs["xim"].Id(p.Userid).Update(p)

	return
}
