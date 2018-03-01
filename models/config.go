package models

import (
	"sync"

	"github.com/silenceper/wechat/cache"
)

// WechatConfig 微信配置
type WechatConfig struct {
	Appid     string
	Appsecret string
	Mutex     *sync.RWMutex
}

// RedisConfig redis配置
type RedisConfig struct {
	Host        string
	Password    string
	Database    int
	MaxIdle     int
	MaxActive   int
	IdleTimeout int32 //second
}
type AppConfigs struct {
	Wechat []WechatConfig
	Redis  RedisConfig
}

// Appconfig 所有微信配置
var Appconfig AppConfigs

// Redisconfig redis配置
var RedisCache *cache.Redis
