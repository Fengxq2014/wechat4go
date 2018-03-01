package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Fengxq2014/wechat4go/models"
	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/oauth"

	"github.com/bitly/go-simplejson"
	"github.com/silenceper/wechat/template"
)

func getConfig(appID string) (config *models.WechatConfig, err error) {
	for _, wechat := range models.Appconfig.Wechat {
		if wechat.Appid == appID {
			return &wechat, nil
		}
	}
	return nil, errors.New("未找到该appid")
}

// GetConfig 获取openid及其他信息
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
		Cache:     models.RedisCache,
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
	result.Unionid = user.UnionID
	if strings.Contains(resToken.Scope, "snsapi_base") == true {
		userInfo, err := wc.GetOauth().GetUserInfo(resToken.AccessToken, resToken.OpenID)
		if err != nil {
			c.Error(errors.New("获取用户信息失败" + err.Error()))
			return
		}
		result.Subscribe = userInfo.Subscribe
		result.Unionid = userInfo.Unionid
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

// GetAccessToken 获取AccessToken
func GetAccessToken(c *gin.Context) {
	type param struct {
		Appid string `form:"appid" json:"appid" binding:"required"`
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

	token := models.Token{
		AppId:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Mutex:     config.Mutex,
	}
	resToken, err := token.GetToken()
	if err != nil {
		c.Error(errors.New("获取token失败" + err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Result{Data: resToken})
}

// SendText 推送客服文本消息
func SendText(c *gin.Context) {
	type param struct {
		Appid   string `from:"appid" json:"appid" binding:"required"`
		Openid  string `from:"openid" json:"openid" binding:"required"`
		Context string `from:"context" json:"context" binding:"required"`
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

	token := models.Token{
		AppId:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Mutex:     config.Mutex,
	}
	resToken, err := token.GetToken()
	if err != nil {
		c.Error(errors.New("获取token失败" + err.Error()))
		return
	}
	sends := models.SendText{
		AccessToken: resToken,
		Jsons: `{
			"touser": "` + queryStr.Openid + `",
			"msgtype":"text",
			"text": {
				"content":"` + queryStr.Context + `"
			}}`,
	}
	err = sends.Send()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, models.Result{Res: 0, Msg: "成功"})
}

// SendImage 发送图片消息
func SendImage(c *gin.Context) {
	type param struct {
		Appid   string `from:"appid" json:"appid" binding:"required"`
		Openid  string `from:"openid" json:"openid" binding:"required"`
		Mediaid string `from:"mediaid" json:"mediaid" binding:"required"`
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
	token := models.Token{
		AppId:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Mutex:     config.Mutex,
	}
	resToken, err := token.GetToken()
	if err != nil {
		c.Error(errors.New("获取token失败" + err.Error()))
		return
	}
	sends := models.SendText{
		AccessToken: resToken,
		Jsons: `{
			"touser": "` + queryStr.Openid + `",
			"msgtype":"image",
			"text": {
				"media_id":"` + queryStr.Mediaid + `"
			}}`,
	}
	err = sends.Send()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, models.Result{Res: 0, Msg: "成功"})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	type param struct {
		Appid  string `from:"appid" json:"appid" binding:"required"`
		Openid string `from:"openid" json:"openid" binding:"required"`
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
	token := models.Token{
		AppId:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Mutex:     config.Mutex,
	}
	resToken, err := token.GetToken()
	if err != nil {
		c.Error(errors.New("获取token失败" + err.Error()))
		return
	}

	wcConfig := &wechat.Config{
		AppID:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Cache:     models.RedisCache,
		Token:     resToken,
	}
	wc := wechat.NewWechat(wcConfig)

	userInfo, err := wc.GetUser().GetUserInfo(queryStr.Openid)
	if err != nil {
		c.Error(errors.New("获取用户信息失败" + err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Result{Res: 0, Data: userInfo})
}

// GetUserInfoOAuth 用户授权获取到用户信息
func GetUserInfoOAuth(c *gin.Context) {
	type param struct {
		Appid  string `from:"appid" json:"appid" binding:"required"`
		Openid string `from:"openid" json:"openid" binding:"required"`
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
	token := models.Token{
		AppId:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Mutex:     config.Mutex,
	}
	resToken, err := token.GetToken()
	if err != nil {
		c.Error(errors.New("获取token失败" + err.Error()))
		return
	}

	user := oauth.Oauth{}
	userInfo, err := user.GetUserInfo(resToken, queryStr.Openid)
	if err != nil {
		c.Error(errors.New("获取用户信息失败" + err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Result{Res: 0, Data: userInfo})
}

// StartKfsession 在线客服
func StartKfsession(c *gin.Context) {
	type param struct {
		Appid  string `from:"appid" json:"appid" binding:"required"`
		Openid string `from:"openid" json:"openid" binding:"required"`
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
	token := models.Token{
		AppId:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Mutex:     config.Mutex,
	}
	resToken, err := token.GetToken()
	if err != nil {
		c.Error(errors.New("获取token失败" + err.Error()))
		return
	}

	send := models.SendText{
		AccessToken: resToken,
	}

	str := send.Get()
	json, err := simplejson.NewJson(str)
	if err != nil {
		c.Error(errors.New("获取客服列表失败" + err.Error()))
		return
	}
	body, err := json.Get("kf_online_list").Array()
	if err != nil {
		c.Error(errors.New("获取客服列表失败" + err.Error()))
		return
	}
	for _, di := range body {
		newdi, _ := di.(map[string]interface{})
		if newdi["status"] == 1 {
			id := newdi["kf_id"]
			sends := models.SendText{
				AccessToken: resToken,
				Jsons: `{
					"touser": "` + queryStr.Openid + `",
					"msgtype":"text",
					"text": {
						"content":"` + "您好，客服" + id.(string) + "为您服务" + `"
					}}`,
			}
			err = sends.Send()
			if err != nil {
				c.Error(err)
				return
			}
		} else {
			sends := models.SendText{
				AccessToken: resToken,
				Jsons: `{
					"touser": "` + queryStr.Openid + `",
					"msgtype":"text",
					"text": {
						"content":"` + "客服正忙，请稍后再试" + `"
					}}`,
			}
			err = sends.Send()
			if err != nil {
				c.Error(err)
				return
			}
		}
	}
	c.JSON(http.StatusOK, models.Result{Res: 0, Data: ""})
}

// SendTemplateMsg 推送模板消息
func SendTemplateMsg(c *gin.Context) {
	type param struct {
		Appid      string `json:"appid" from:"appid" binding:"required"`
		Openid     string `json:"openid" from:"openid" binding:"required"`
		Url        string `json:"url" from:"url"`
		TemplateId string `json:"templateId" from:"templateId" binding:"required"`
		First      string `json:"first" from:"first"`
		Keyword1   string `json:"keyword1" from:"keyword1"`
		Keyword2   string `json:"keyword2" from:"keyword2"`
		Keyword3   string `json:"keyword3" from:"keyword3"`
		Keyword4   string `json:"keyword4" from:"keyword4"`
		Keyword5   string `json:"keyword5" from:"keyword5"`
		Remark     string `json:"remark" from:"remark"`
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
		AppID:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Cache:     models.RedisCache,
	}
	wc := wechat.NewWechat(wcConfig)

	json := map[string]*template.DataItem{}
	json["Value"].Value = `{
		"first": {
			"value":"您好，` + queryStr.First + `，您有一份完整测评报告已生成。\r\n",
			"color":""
		},
		"keyword1":{
			"value":"` + queryStr.Keyword1 + `",
			"color":""
		},
		"keyword2": {
			"value":"` + queryStr.Keyword2 + `",
			"color":""
		},
		"keyword3": {
			"value":"` + queryStr.Keyword3 + `\r\n",
			"color":""
		},
		"keyword4": {
			"value":"` + queryStr.Keyword4 + `\r\n",
			"color":""
		},
		"keyword5": {
			"value":"` + queryStr.Keyword5 + `\r\n",
			"color":""
		},
		"remark":{
			"value":"` + queryStr.Remark + `",
			"color":"#89bd41"
		}}`
	json["Color"].Color = "#89bd41"

	msg := template.Message{
		ToUser:     queryStr.Openid,
		TemplateID: queryStr.TemplateId,
		URL:        queryStr.Url,
		Data:       json,
	}

	id, err := wc.GetTemplate().Send(&msg)
	if err != nil {
		c.Error(errors.New("发送模板消息失败" + err.Error()))
		return
	}
	if id != 0 {
		c.Error(errors.New("发送模板消息失败" + err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Result{Res: 0, Data: ""})
}

// DownloadTempMedia 获取临时素材
func DownloadTempMedia(c *gin.Context) {
	type param struct {
		Appid   string `from:"appid" json:"appid" binding:"required"`
		MediaId string `from:"mediaId" json:"mediaId" binding:"required"`
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
		AppID:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Cache:     models.RedisCache,
	}
	wc := wechat.NewWechat(wcConfig)

	url, err := wc.GetMaterial().GetMediaURL(queryStr.MediaId)

	if err != nil {
		c.Error(errors.New("获取临时素材URL失败" + err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.Result{Res: 0, Data: url})
}

// GetJSSDKConfig  获取jssdk配置信息
func GetJSSDKConfig(c *gin.Context) {
	type param struct {
		Appid string `from:"appid" json:"appid"`
		Url   string `from:"url" json:"url"`
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
		AppID:     queryStr.Appid,
		AppSecret: config.Appsecret,
		Cache:     models.RedisCache,
	}
	wc := wechat.NewWechat(wcConfig)

	con, err := wc.GetJs().GetConfig(queryStr.Url)
	if err != nil {
		c.Error(errors.New("获取jssdk配置信息失败" + err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Result{Res: 0, Data: con})
}
