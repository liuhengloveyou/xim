/*
* 好友逻辑服务
 */

package face

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/liuhengloveyou/xim/service"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

func init() {
	http.HandleFunc("/friends/list", friendsList)
}

func friendsList(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "GET" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "只支持GET请求")
		return
	}

	sess, auth := authFilter(w, r)
	if auth == false {
		log.Errorln("friendsList ERR: 末登录用户.")
		gocommon.HttpErr(w, http.StatusForbidden, "末登录用户.")
		return
	}
	userid := sess.Get("user").(*service.User).Userid

	ver := strings.TrimSpace(r.FormValue("v"))
	if ver == "" {
		ver = "0"
	}

	iver, e := strconv.ParseUint(ver, 10, 64)
	if e != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, e.Error())
		return
	}

	result, e := service.FriendList(userid, uint(iver))
	if e != nil {
		log.Errorln("friendsList ERR:", e.Error())
		gocommon.HttpErr(w, http.StatusInternalServerError, "数据库服务错误.")
		return
	}

	log.Infoln("friendsList:", iver, result)
	if _, e = w.Write([]byte(result)); e != nil {
		log.Exitln(e)
	}
}
