package httpclient_test

import (
	"testing"
	"github.com/ydbt/devtool/httpclient"
)

func TestHttpGet(t *testing.T) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json;charset=UTF-8"
	header["abcd"] = "1234"
	hc := httpclient.NewHttpClient(header)
	body, err := hc.Get("http://www.baidu.com/", map[string]interface{}{"a": 100, "0": 3.14})
	if err != nil {
		t.Error(err)
	}
	t.Log(string(body))
}

func TestHttpPost(t *testing.T) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json;charset=UTF-8"
	header["abcd"] = "1234"
	hc := httpclient.NewHttpClient(header)
	body, err := hc.Post("http://www.baidu.com/", nil, "hello world")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(body))
}
