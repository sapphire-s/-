package main

import (
	"scutbot.cn/uniauth/bootstrap"
	"scutbot.cn/uniauth/log"
	"scutbot.cn/uniauth/router"
)

func main() {
	log.Init()
	bootstrap.InitSetting()
	router.InitRouter()

}
