package face

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"
	"github.com/liuhengloveyou/xim/service"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

func init() {
	http.HandleFunc("/recv", recvMessage)
	http.HandleFunc("/send", sendMessage)
}

func HttpAccess() {
	//http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("/Users/liuheng/go/src/github.com/liuhengloveyou/xim-ionic/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome you!"))
		log.Infoln("RequestURI:", r.RequestURI)
	})

	s := &http.Server{
		Addr:           fmt.Sprintf("%v:%v", common.AccessConf.Addr, common.AccessConf.Port),
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("HTTP IM GO... %v:%v\n", common.AccessConf.Addr, common.AccessConf.Port)
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func optionsFilter(w http.ResponseWriter, r *http.Request) {
	return

	w.Header().Set("Access-Control-Allow-Origin", "http://web.xim.com:9000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "X-API, X-REQUEST-ID, X-API-TRANSACTION, X-API-TRANSACTION-TIMEOUT, X-RANGE, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)

	return
}

func authFilter(w http.ResponseWriter, r *http.Request) (sess session.SessionStore, auth bool) {
	token := strings.TrimSpace(r.Header.Get("TOKEN"))
	if token == "" {
		if cookie, e := r.Cookie(common.AccessConf.Session.CookieName); e == nil {
			if cookie != nil {
				token = cookie.Value
			}
		}
	}
	if token == "" {
		log.Errorln("token nil")
		return nil, false
	}

	sess, err := session.GetSessionById(token)
	if err != nil {
		log.Warningln("session ERR:", err.Error())
		return nil, false
	}
	log.Infoln("auth:", token, sess)

	if sess.Get("user") != nil {
		return sess, true
	}

	// passport auth.
	log.Errorln("passport auth:", sess)
	info, err := common.Passport.UserAuth(token)
	if err != nil {
		log.Errorln("passport auth ERR:", token, err.Error())
		return nil, false
	}

	user := &service.User{}
	if err := json.Unmarshal(info, user); err != nil {
		log.Errorln("passport response ERR:", token, string(info))
		return nil, false
	}

	sess.Set("user", user)              // 用户信息
	sess.Set("sync", time.Now().Unix()) // 同步时间
	log.Errorln("session from passport:", sess)
	return sess, true
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "只支持POST请求.")
		return
	}

	sess, auth := authFilter(w, r)
	if auth == false {
		log.Errorln("send ERR: 末登录用户.")
		gocommon.HttpErr(w, http.StatusForbidden, "末登录用户.")
		return
	}

	api := r.Header.Get("X-API")
	log.Infoln("sendMessage X-API:", api)

	switch api {
	case common.LOGIC_TEMPGROUP:
		if _, e := tgroup(r, "send"); e != nil {
			log.Errorln("sendMessage tgroup ERR:", e.Error())
			gocommon.HttpErr(w, http.StatusInternalServerError, "临时讨论组系统错误.")
			return
		}
	default:
		body, e := ioutil.ReadAll(r.Body)
		if e != nil {
			log.Errorln("get request.body ERR:", e.Error())
			gocommon.HttpErr(w, http.StatusBadRequest, "请求错误.")
			return

		}

		if e := service.SendMessage(sess, body); e != nil {
			log.Errorln("sendMessage ERR:", e.Error())
			gocommon.HttpErr(w, http.StatusInternalServerError, "系统错误.")
			return
		}
	}

	gocommon.HttpErr(w, http.StatusOK, "OK")

	return
}

func recvMessage(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	sess, auth := authFilter(w, r)
	if auth == false {
		log.Errorln("recv ERR: 末登录用户.")
		gocommon.HttpErr(w, http.StatusForbidden, "末登录用户.")
		return
	}

	api := r.Header.Get("X-API")
	log.Infoln("recvMessage X-API:", api)

	var (
		info *service.UserSession
		e    error
	)

	switch api {
	case common.LOGIC_TEMPGROUP:
		if info, e = tgroup(r, "recv"); e != nil {
			gocommon.HttpErr(w, http.StatusInternalServerError, "临时讨论组系统错误.")
			return
		}
	default:
		if info, e = service.StateUpdate(sess); e != nil {
			log.Errorln("StateUpdate ERR:", e)
			gocommon.HttpErr(w, http.StatusInternalServerError, "系统错误.")
			return
		}
	}

	if info == nil {
		gocommon.HttpErr(w, http.StatusInternalServerError, "系统内部错误.")
		return
	}
	if info.ID == "" {
		gocommon.HttpErr(w, http.StatusInternalServerError, "系统内部错误.")
		return
	}

	ctx := ""
	select {
	case ctx = <-info.MsgChan:
	case <-time.After(1 * time.Minute):
		ctx = "TIMEOUT" // 长连接每1分钟断开一次, 没有心跳
	case <-w.(http.CloseNotifier).CloseNotify():
		log.Warningln("client closed:", api, sess.Id(""), sess.Get("user"))
	}

	log.Infoln("Recv OK:", sess.Id(""), sess.Get("user"), ctx)
	w.Write([]byte(ctx))

	return
}

// 临时讨论组
func tgroup(r *http.Request, logic string) (user *service.UserSession, e error) {
	if "recv" == logic {
		userid, groupid := r.FormValue("uid"), r.FormValue("gid")
		if userid == "" || groupid == "" {
			return nil, fmt.Errorf("末知的用户或组.")
		}
		log.Infoln("tgroup: ", userid, groupid)

		if user, e = service.TGroutRecv(userid, groupid); e != nil {
			return nil, e
		}

		return user, nil
	} else if "send" == logic {
		userid, groupid := r.FormValue("uid"), r.FormValue("gid")
		if userid == "" || groupid == "" {
			return nil, fmt.Errorf("末知的用户或组.")
		}
		bm, e := ioutil.ReadAll(r.Body)
		if e != nil {
			return nil, e
		}
		log.Infoln("tgroup: ", userid, groupid, string(bm))

		if e = service.TGroutSend(userid, groupid, string(bm)); e != nil {
			return nil, e
		}
	}

	return nil, nil
}
