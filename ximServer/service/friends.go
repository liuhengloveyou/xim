package service

import (
	"github.com/liuhengloveyou/xim/dao"
)

func FriendList(userid string, version uint) (result string, e error) {
	one := &dao.Friends{Userid: userid, Version: int(version)}
	if e = one.GetOneByVersion(); e != nil {
		return "", e
	}
	if one.Friends != nil {
		return *one.Friends, nil
	}

	return "", nil
}
