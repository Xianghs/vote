package models

import (
	"log"
	"vote/vote/app/db"
	"vote/vote/app/util"
)

//建表
func Init_dbtable() {
	var mydb = db.MyDB.AutoMigrate(
		&User{},
		&Vote{},
		&VoteItem{},
		&VoteRecord{},
		&Admin{},
		&Notice{},
		&SystemSet{},
	)
	if err := mydb.Error; err != nil {
		log.Fatal("初始化数据库表失败：", err)
	}

	//创建一个超级管理员帐户
	if err := InitAdmin(); err != nil {
		log.Fatal("超级管理员初始化失败:", err)
	}

	if err := InitSystem(); err != nil {
		log.Fatal("系统初始化失败：", err)
	}
}

//创建一个超级管理员帐户
func InitAdmin() error {
	//先查询是否已经存在超级管理员
	count := 0
	if err := db.MyDB.Model(&Admin{}).Where("level=1").Count(&count).Error; err != nil {
		return err
	}

	//如果不存在就创建
	if count == 0 {
		pwd, err := util.AesEncrypt("123456")
		if err != nil {
			return err
		}
		adm := Admin{
			ID:        util.CreateRandomString(20),
			AdminName: "admin",
			Level:     1,
			Password:  pwd,
		}
		if err := db.MyDB.Model(&Admin{}).Create(&adm).Error; err != nil {
			return err
		}
	}
	return nil
}

//初始化系统设置
func InitSystem() error {
	//查询记录是否存在
	count := 0
	if err := db.MyDB.Model(&SystemSet{}).Count(&count).Error; err != nil {
		return err
	}
	//不存在就创建
	if count == 0 {
		sys := SystemSet{
			ID:             1,
			TotalHousehold: 999,
			TotalArea:      999,
		}

		if err := db.MyDB.Model(&SystemSet{}).Create(&sys).Error; err != nil {
			return err
		}
	}

	return nil
}
