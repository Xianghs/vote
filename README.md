# vote
2020年本科毕业设计后端部分代码

#直接使用需要补充如下配置信息
``` ini vote/app/config/config.ini
;==================== Mysql
[mysql]
Type         	= mysql
Host            = 数据库IP
Port            = 3306
User            = vote
Password        = 数据库密码
Database        = vote
Charset      	= utf8mb4
MaxIdleConnects = 5
MaxOpenConnects = 10

;==================== Redis
[redis]
Host        = Redis IP
Port        = 6379
Password    = redis密码
MaxIdle     = 30
MaxActive   = 30
IdleTimeout = 200

;==================== 对象存储，存储空间权限开放，url用http
[ali_oss]
Endpoint        = 
AccessKeyId     = 
AccessKeySecret = 
BucketName      = 

```
