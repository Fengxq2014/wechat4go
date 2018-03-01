package route

import (
	"github.com/Fengxq2014/wechat4go/api"
	"github.com/gin-gonic/gin"
)

var DB = make(map[string]string)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(handleErrors)
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})
	r.POST("/Weixin/GetConfig", api.GetConfig)
	r.POST("/Weixin/GetAccessToken", api.GetAccessToken)
	r.POST("/Weixin/SendText", api.SendText)
	r.POST("/Weixin/SendImage", api.SendImage)
	r.POST("/Weixin/GetUserInfo", api.GetUserInfo)
	r.POST("/Weixin/GetUserInfoOAuth", api.GetUserInfoOAuth)
	r.POST("/Weixin/StartKfsession", api.StartKfsession)
	r.POST("/Weixin/SendTemplateMsg", api.SendTemplateMsg)
	r.POST("/Weixin/DownloadTempMedia", api.DownloadTempMedia)
	r.POST("/Weixin/GetJSSDKConfig", api.GetJSSDKConfig)

	return r
}

func handleErrors(c *gin.Context) {
	c.Next()
	errorToPrint := c.Errors.Last()
	if errorToPrint != nil {
		c.JSON(200, gin.H{
			"res":  1,
			"msg":  errorToPrint.Error(),
			"data": nil,
		})
	}
}
