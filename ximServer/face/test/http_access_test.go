package access_test

import (
	"net/http"
	"testing"

	gocommon "github.com/liuhengloveyou/go-common"
)

const url = "http://127.0.0.1:6060"

var user1 = `
{
	"cellphone": "18510511015",
	"password": "123456"
}
`
var user2 = `
{
	"email": "liuhengloveyou@gmail.com",
	"password": "123456"
}
`
var cookies []*http.Cookie

func TestUserRegister(t *testing.T) {
	statuCode, responseCookies, responseBody, err := gocommon.PostRest(url+"/user/add", []byte(user1), nil, nil)
	t.Log(statuCode)
	t.Log(responseCookies)
	t.Log(string(responseBody))
	t.Log(err)

	statuCode, responseCookies, responseBody, err = gocommon.PostRest(url+"/user/add", []byte(user2), nil, nil)
	t.Log(statuCode)
	t.Log(responseCookies)
	t.Log(string(responseBody))
	t.Log(err)
}

func TestUserLogin(t *testing.T) {
	statuCode, responseCookies, responseBody, err := gocommon.PostRest(url+"/user/login", []byte(user1), nil, nil)
	cookies = responseCookies
	t.Log(statuCode)
	t.Log(responseCookies)
	t.Log(string(responseBody))
	t.Log(err)

}

func TestUserAuth(t *testing.T) {
	statuCode, responseCookies, responseBody, err := gocommon.PostRest(url+"/user/auth", nil, cookies, nil)
	t.Log(cookies)
	t.Log(statuCode)
	t.Log(responseCookies)
	t.Log(string(responseBody))
	t.Log(err)

}
