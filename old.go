package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Fengxq2014/wechat4go/models"
	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat"
)

func getConfig(appID string) (config *WechatConfig, err error) {
	for _, wechat := range Appconfig.Wechat {
		if wechat.Appid == appID {
			return &wechat, nil
		}
	}
	return nil, errors.New("未找到改appid")
}

func GetConfig(c *gin.Context) {
	type param struct {
		Appid string `form:"appid" json:"appid" binding:"required"`
		Code  string `form:"code" json:"code" binding:"required"`
		JsURL string `form:"jsUrl" json:"jsUrl"`
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.JSON) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	config, err := getConfig(queryStr.Appid)
	if err != nil {
		c.Error(err)
		return
	}
	wcConfig := &wechat.Config{
		AppID:     config.Appid,
		AppSecret: config.Appsecret,
		Cache:     RedisCache,
	}
	result := models.GetConfigRespose{}
	wc := wechat.NewWechat(wcConfig)
	resToken, err := wc.GetOauth().GetUserAccessToken(queryStr.Code)
	if err != nil {
		c.Error(errors.New("获取openid失败" + err.Error()))
		return
	}
	result.Openid = resToken.OpenID
	user, err := wc.GetUser().GetUserInfo(resToken.OpenID)
	if err != nil {
		c.Error(errors.New("GetUserInfo失败" + err.Error()))
		return
	}
	result.Subscribe = user.Subscribe
	result.Unionid = user.UnionID
	if strings.Contains(resToken.Scope, "snsapi_base") == true {
		userInfo, err := wc.GetOauth().GetUserInfo(resToken.AccessToken, resToken.OpenID)
		if err != nil {
			c.Error(errors.New("获取用户信息失败" + err.Error()))
			return
		}
	}
	if queryStr.JsURL != "" {
		jsconfig, err := wc.GetJs().GetConfig(queryStr.JsURL)
		if err != nil {
			c.Error(errors.New("获取jsconfig失败" + err.Error()))
			return
		}
		result.Nonce = jsconfig.NonceStr
		result.Signature = jsconfig.Signature
		result.Timestamp = jsconfig.Timestamp
	}
	c.JSON(http.StatusOK, models.Result{Data: &result})
}
