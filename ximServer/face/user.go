/*
* 用户信息逻辑服务
 */

package face

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/xim/service"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/validator"
)

func init() {
	http.HandleFunc("/user/add", UserAdd)
}

func UserAdd(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "POST" {
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, "只支持POST请求")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		return

	}
	log.Infoln(string(body))

	user := &service.User{}
	if err = json.Unmarshal(body, user); err != nil {
		log.Errorln("json.Unmarshal(body, user) ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = validator.Validate(user); err != nil {
		log.Errorln("validator ERR: ", err)
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		return
	}

	if user.Cellphone == "" && user.Email == "" {
		log.Errorln("ERR: 用户手机号和邮箱地址同时为空.")
		gocommon.HttpErr(w, http.StatusBadRequest, "用户手机号和邮箱地址同时为空.")
		return
	}
	if user.Password == "" {
		log.Errorln("ERR: 用户密码为空.")
		gocommon.HttpErr(w, http.StatusBadRequest, "用户密码为空.")
		return
	}

	if err = user.AddUser(); err != nil {
		log.Errorln("AddUser ERR: ", err)
		gocommon.HttpErr(w, http.StatusInternalServerError, "系统错误,请联系管理员.")
		return

	}

	log.Infoln("add user:", user)
	fmt.Fprintf(w, "{\"userid\":\"%s\"}", user.Userid)
	return
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	optionsFilter(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gocommon.HttpErr(w, http.StatusBadRequest, err.Error())
		log.Errorln("ioutil.ReadAll(r.Body) ERR: ", err)
		return
	}
	log.Infoln(r.RequestURI, string(body))

	/*
		stat, cookies, response, e := Passport.UserAdd(r.RequestURI, body, r.Cookies())
		if e != nil {
			gocommon.HttpErr(w, http.StatusInternalServerError, e.Error())
			log.Errorln("call passport ERR: ", err)
			return
		}
		fmt.Println(stat, string(response), e)

		if cookies != nil {
			for _, cookie := range cookies {
				http.SetCookie(w, cookie)
			}
		}

		gocommon.HttpErr(w, stat, string(response))
	*/

	return
}
