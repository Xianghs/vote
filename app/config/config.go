package config

import (
	"github.com/gin-contrib/cors"
	"gopkg.in/ini.v1"
	"log"
	"time"
)

type DevConfig struct {
	Env     string
	LogPath string
	Host    string
	Port    string
}

type MysqlConfig struct {
	Type            string
	User            string
	Password        string
	Host            string
	Port            int
	Charset         string
	Database        string
	MaxIdleConnects int
	MaxOpenConnects int
}

type RedisConfig struct {
	Host        string
	Port        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type LimiterConfig struct {
	Limit float64
	Burst int
}

type AESConfig struct {
	AesKey string
}

type AliOSSConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
}

var Dev DevConfig
var MySql MysqlConfig
var Redis RedisConfig
var Limiter LimiterConfig

var Cors cors.Config
var AES AESConfig
var AliOSS AliOSSConfig

func SetUpConfig(configPath string) {
	//configPath := "./config/config.ini"
	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatal("解析配置文件失败：", err)
	}
	getSection(cfg, "dev", &Dev)
	getSection(cfg, "mysql", &MySql)
	getSection(cfg, "redis", &Redis)
	getSection(cfg, "limiter", &Limiter)
	Cors = cors.DefaultConfig()
	getSection(cfg, "cors", &Cors)
	getSection(cfg, "aes", &AES)
	getSection(cfg, "ali_oss", &AliOSS)
}

func getSection(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatal("读取配置文件section失败", err)
	}
}
