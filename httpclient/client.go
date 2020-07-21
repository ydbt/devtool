package httpclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpClient struct {
	client    *http.Client
	ctxHeader map[string]string
}

func NewHttpClient(header map[string]string) *HttpClient {
	hc := new(HttpClient)
	hc.client = new(http.Client)
	hc.ctxHeader = make(map[string]string)
	for k, v := range header {
		hc.ctxHeader[k] = v
	}
	return hc
}

func (hc *HttpClient) Get(uri string, param map[string]interface{}) ([]byte, error) {
	pv := url.Values{}
	for k, v := range param {
		pv.Set(k, fmt.Sprintf("%v", v))
	}
	req, err := http.NewRequest("GET", uri+pv.Encode(), nil)
	res, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

func (hc *HttpClient) Post(uri string, param map[string]interface{}, reqBody string) ([]byte, error) {
	pv := url.Values{}
	for k, v := range param {
		pv.Set(k, fmt.Sprintf("%v", v))
	}
	req, err := http.NewRequest("POST", uri+pv.Encode(), strings.NewReader(reqBody))
	for k, v := range hc.ctxHeader {
		req.Header.Add(k, v)
	}
	res, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	replyBody, err := ioutil.ReadAll(res.Body)
	return replyBody, err
}
