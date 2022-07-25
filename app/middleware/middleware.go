package middleware

import (
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
	"vote/vote/app/config"
	"vote/vote/app/errCode"
	"vote/vote/app/util"
	"vote/vote/app/util/resp"
)

//日志中间件
func Logger() gin.HandlerFunc {

	logClient := log.New()
	var logPath = config.Dev.LogPath // 日志打印到指定的目录

	fileName := path.Join(logPath, "api.log")
	//禁止logrus的输出
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//如果文件不存在就创建
	if err != nil && os.IsNotExist(err) {
		src, _ = os.Create(fileName)
	}

	// 设置日志输出的路径
	logClient.Out = src
	logClient.SetLevel(log.DebugLevel)
	//apiLogPath := "gin-api.log"
	logWriter, err := rotatelogs.New(
		fileName+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(fileName),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	writeMap := lfshook.WriterMap{
		log.InfoLevel:  logWriter,
		log.FatalLevel: logWriter,
		log.DebugLevel: logWriter, // 为不同级别设置不同的输出目的
		log.WarnLevel:  logWriter,
		log.ErrorLevel: logWriter,
		log.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &log.JSONFormatter{})
	logClient.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		// 这里是指定日志打印出来的格式。分别是状态码，执行时间,请求ip,请求方法,请求路由(等下我会截图)
		logClient.Infof("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}

//访问频率限制
var limiter *util.IPRateLimiter
func LimitMiddleware() gin.HandlerFunc {
	limiter = util.NewIPRateLimiter(config.Limiter.Limit, config.Limiter.Burst)
	return func(c *gin.Context) {
		var resp = resp.Resp{c}
		limiter := limiter.GetLimiter(resp.Request.RemoteAddr)
		if !limiter.Allow() {
			resp.ResponseErr(errCode.TooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}
