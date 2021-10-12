package httpkit

import (
	"bytes"
	"crypto/tls"
	"github.com/myleo1/go-core-kit/library/jsonkit"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

var client *http.Client

func init() {
	// 忽略证书校验 todo
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar
}

// Req 填写FormData、JsonData时可缺省contentType
type Req struct {
	Method      string
	Url         string
	Header      map[string]string
	ContentType string
	FormData    map[string]string
	JsonData    interface{}
	BinaryData  []byte
	Timeout     int // seconds
}

type Resp struct {
	Payload *http.Response
}

const ContentTypeForm = "application/x-www-form-urlencoded; charset=utf-8"
const ContentTypeJSON = "application/json; charset=utf-8"

func Request(reqBean Req) (*Resp, error) {
	if reqBean.Method == "" {
		reqBean.Method = http.MethodPost
	}
	var req *http.Request
	var err error
	if reqBean.BinaryData != nil {
		req, err = http.NewRequest(reqBean.Method, reqBean.Url, bytes.NewBuffer(reqBean.BinaryData))
	} else if reqBean.JsonData != nil {
		req, err = http.NewRequest(reqBean.Method, reqBean.Url, bytes.NewBuffer([]byte(jsonkit.ToString(reqBean.JsonData))))
	} else {
		data := make(url.Values)
		for key, val := range reqBean.FormData {
			data.Add(key, val)
		}
		req, err = http.NewRequest(reqBean.Method, reqBean.Url, strings.NewReader(data.Encode()))
	}
	if err != nil {
		return nil, err
	}
	if reqBean.ContentType == "" {
		if reqBean.JsonData != nil {
			req.Header.Set("Content-Type", ContentTypeJSON)
		} else {
			req.Header.Set("Content-Type", ContentTypeForm)
		}
	} else {
		req.Header.Set("Content-Type", reqBean.ContentType)
	}
	for key, val := range reqBean.Header {
		req.Header.Set(key, val)
	}
	if reqBean.Timeout > 0 {
		client.Timeout = time.Duration(reqBean.Timeout) * time.Second
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	res := &Resp{
		Payload: resp,
	}
	return res, nil
}

func (th *Resp) RespBody2Str() (string, error) {
	defer th.Payload.Body.Close()
	body, err := io.ReadAll(th.Payload.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (th *Resp) StatusCode() int {
	return th.Payload.StatusCode
}

func (th *Resp) GetQueryParam(key string) string {
	return th.Payload.Request.URL.Query().Get(key)
}
