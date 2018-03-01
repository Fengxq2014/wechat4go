package models

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/bitly/go-simplejson"
)

const (
	incompleteURL = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="
	kefuUrl       = "https://api.weixin.qq.com/cgi-bin/customservice/getonlinekflist?access_token="
)

type SendText struct {
	AccessToken string `from:"accesstoken" json:"accesstoken"`
	Jsons       string `from:"jsons" json:"jsons"`
}

// Send post方法
func (c *SendText) Send() (err error) {
	b := []byte(c.Jsons)
	body := bytes.NewBuffer(b)

	res, err := http.Post(incompleteURL+c.AccessToken, "application/json;charset=utf-8", body)

	result, err := ioutil.ReadAll(res.Body)
	result, _ = GbkToUtf8(result)
	res.Body.Close()

	js, err := simplejson.NewJson(result)
	if err != nil {
		return
	}
	if js.Get("errcode").MustInt() == 0 {
		return nil
	}

	return errors.New(string(result))
}

// Get Get方法
func (c *SendText) Get() (ret []byte) {
	res, err := http.Get(kefuUrl + c.AccessToken)
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return result
}

// GbkToUtf8 方法
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
