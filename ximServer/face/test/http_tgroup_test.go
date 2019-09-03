package access_test

import (
	"testing"

	gocommon "github.com/liuhengloveyou/go-common"
)

const url = "http://127.0.0.1:6060"

func TestTGroupLogin(t *testing.T) {
	statuCode, responseCookies, responseBody, err := gocommon.PostRest(url+"/recv?uid=123&gid=g123", nil, nil, &map[string]string{"X-API": "tgroup"})
	t.Log(statuCode)
	t.Log(responseCookies)
	t.Log(string(responseBody))
	t.Log(err)
}

func TestTGroupPush(t *testing.T) {
	statuCode, responseCookies, responseBody, err := gocommon.PostRest(url+"/send?uid=123&gid=g123", nil, nil, &map[string]string{"X-API": "tgroup"})
	t.Log(statuCode)
	t.Log(responseCookies)
	t.Log(string(responseBody))
	t.Log(err)
}
