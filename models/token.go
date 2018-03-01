package models

import (
	"sync"

	"github.com/silenceper/wechat/context"
)

type Token struct {
	AppId     string `from:"appid" json:"appid"`
	AppSecret string `from:"appsecret" json:"appsecret"`
	Mutex     *sync.RWMutex
}

// 获取 AccessToken
func (token *Token) GetToken() (tokens string, err error) {
	ctx := context.Context{
		AppID:     token.AppId,
		AppSecret: token.AppSecret,
		Cache:     RedisCache,
	}

	ctx.SetAccessTokenLock(token.Mutex)
	tokens, err = ctx.GetAccessToken()
	if err != nil {
		return
	}
	return
}
