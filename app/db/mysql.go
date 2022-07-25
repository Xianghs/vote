package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"vote/vote/app/config"
)

//数据库连接
var MyDB *gorm.DB

func InitMysql(cfg config.MysqlConfig) {
	var err error
	var url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host,
		cfg.Port, cfg.Database, cfg.Charset)

	if MyDB, err = gorm.Open(cfg.Type, url); err != nil {
		log.Fatal("Mysql连接失败：err", err)
		return
	}

	if config.Dev.Env == "development" {
		MyDB.LogMode(true)
	}

	MyDB.DB().SetMaxIdleConns(config.MySql.MaxIdleConnects)
	MyDB.DB().SetMaxOpenConns(config.MySql.MaxOpenConnects)

	MyDB.SingularTable(true)

}
