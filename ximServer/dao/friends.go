package dao

import (
	"database/sql"

	"github.com/liuhengloveyou/xim/common"
)

type Friends struct {
	Userid  string  `xorm:"not null pk INT(11)"`
	Friends *string `xorm:"JSON"`
	Version int     `xorm:"INT(11)"`
}

func (p *Friends) Insert() (e error) {
	_, e = common.DBs["xim"].Insert("INSERT INTO friends values(?,?,?)", p.Userid, p.Friends, 1)

	return
}

func (p *Friends) Find() (one []*Friends, e error) {
	//	e = common.DBs["xim"].Query(sqlStr string, args ...interface{})

	return
}

func (p *Friends) GetOne() (has bool, e error) {
	//	has, e = common.DBs["xim"].Get(p)

	return
}

func (p *Friends) GetOneByVersion() (e error) {
	var rst sql.NullString
	if e = common.DBs["xim"].Conn.QueryRow("SELECT * FROM friends WHERE userid=? and version > ?;", p.Userid, p.Version).Scan(&p.Userid, &rst, &p.Version); e != nil {
		if e == sql.ErrNoRows {
			e = nil
		}

		return
	}
	if rst.Valid {
		p.Friends = &rst.String
	}

	return
}
