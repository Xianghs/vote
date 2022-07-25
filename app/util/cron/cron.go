package cron

import (
	log "github.com/sirupsen/logrus"
	"time"
	"vote/vote/app/db"
	"vote/vote/app/models"
)

//定时任务

//启动多个任务
func InitCron() {
	t1 := time.NewTimer(time.Second * 60)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 60)
			vote()
		}
	}
}

func vote() {
	if err := db.MyDB.Model(&models.Vote{}).
		Where("end_time<? and status=?", time.Now(), "进行中").
		Update("status", "已结束").Error; err != nil {
		log.Warn("投票到期调整错误！！！！！！！：", err.Error())
	}
}
