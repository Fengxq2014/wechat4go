package main

import "github.com/Fengxq2014/wechat4go/route"

func main() {
	r := route.SetupRouter()
	r.Run(":3000") // listen and serve on 0.0.0.0:8080
}
