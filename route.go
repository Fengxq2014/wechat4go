package main

import "github.com/gin-gonic/gin"

var DB = make(map[string]string)

func setupRouter() *gin.Engine {
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
