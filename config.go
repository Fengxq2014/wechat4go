package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/silenceper/wechat/cache"

	"github.com/Fengxq2014/wechat4go/models"
)

func init() {
	if _, err := toml.DecodeFile("./conf/cust.toml", &models.Appconfig); err != nil {
		log.Fatalln("config error " + err.Error())
	}
	fmt.Println("config length ", len(models.Appconfig.Wechat))
	opts := &cache.RedisOpts{
		Host:        models.Appconfig.Redis.Host,
		Password:    models.Appconfig.Redis.Password,
		Database:    models.Appconfig.Redis.Database,
		MaxActive:   models.Appconfig.Redis.MaxActive,
		MaxIdle:     models.Appconfig.Redis.MaxIdle,
		IdleTimeout: models.Appconfig.Redis.IdleTimeout,
	}
	models.RedisCache = cache.NewRedis(opts)

	for i, _ := range models.Appconfig.Wechat {
		models.Appconfig.Wechat[i].Mutex = new(sync.RWMutex)
	}
}
