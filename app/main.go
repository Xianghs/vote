package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"vote/vote/app/config"
	"vote/vote/app/db"
	"vote/vote/app/models"
	"vote/vote/app/router"
	"vote/vote/app/util/cron"
)

func main() {
	//初始化配置
	config.SetUpConfig("./config/config.ini")

	//链接数据库
	db.InitMysql(config.MySql)
	db.InitRedis(config.Redis)
	//自动建表
	models.Init_dbtable()

	var address = fmt.Sprintf("%s:%s", config.Dev.Host, config.Dev.Port)
	var app = gin.New()
	if config.Dev.Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	//初始化路由
	router.InitRouter(app)

	//初始化定时任务
	go cron.InitCron()

	if err := app.Run(address); err != nil {
		log.Fatal("服务启动失败：", err)
	}

}
