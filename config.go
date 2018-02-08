package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/silenceper/wechat/cache"
)

// WechatConfig 微信配置
type WechatConfig struct {
	Appid     string
	Appsecret string
}

// RedisConfig redis配置
type redisConfig struct {
	Host        string
	Password    string
	Database    int
	MaxIdle     int
	MaxActive   int
	IdleTimeout int32 //second
}
type appConfigs struct {
	Wechat []WechatConfig
	redis  redisConfig
}

// Appconfig 所有微信配置
var Appconfig appConfigs

// Redisconfig redis配置
var RedisCache *cache.Redis

func init() {
	if _, err := toml.DecodeFile("./conf/cust.toml", &Appconfig); err != nil {
		log.Fatalln("config error " + err.Error())
	}
	fmt.Println("config length ", len(Appconfig.Wechat))
	opts := &cache.RedisOpts{
		Host:        Appconfig.redis.Host,
		Password:    Appconfig.redis.Password,
		Database:    Appconfig.redis.Database,
		MaxActive:   Appconfig.redis.MaxActive,
		MaxIdle:     Appconfig.redis.MaxIdle,
		IdleTimeout: Appconfig.redis.IdleTimeout,
	}
	RedisCache = cache.NewRedis(opts)
}
