package admin_token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"time"
	"vote/vote/app/db"
	"vote/vote/app/errCode"
	"vote/vote/app/models"
	"vote/vote/app/util/resp"
)

type TokenHandler func(resp *resp.Resp, tokenInfo *Claims)

// 指定加密密钥
var jwtSecret = []byte("")

//Claim是一些实体（通常指的用户）的状态和额外的元数据
type Claims struct {
	ID    string `json:"id"`
	Level int    `json:"level"`
	jwt.StandardClaims
}

// 产生token,并保存到redis
func GenerateToken(admin models.Admin) string {
	//设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		ID:    admin.ID,
		Level: admin.Level,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "vote",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		log.Error("生成token失败：", err)
	}
	key := "admin_token+" + admin.ID

	cmd := redis.NewStringCmd("set", key, token, "ex", "2505600")
	if err := db.MyRedis.Process(cmd); err != nil {
		log.Error("保存token错误：", err)
		return ""
	}

	if cmd.Err() != nil {
		log.Error("保存token失败：", cmd.Err().Error())
	}
	return token
}

// 根据传入的token值获取到Claims对象信息，（进而获取其中的ID,用户名和密码）
func ParseToken(token string) (*Claims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

//token验证的中间件
func RequiredToken(handler TokenHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			tokenInfo *Claims
			resp      = resp.NewResp(c)
		)
		tokenStr := c.GetHeader("token")
		tokenInfo, err := ParseToken(tokenStr)
		if err != nil {
			resp.ResponseErr(errCode.TokenInvalid)
			resp.Abort()
			return
		}
		key := "admin_token+" + tokenInfo.ID
		//在redis中查找token
		var str string
		if err := db.MyRedis.Get(key).Scan(&str); err != nil {
			resp.ResponseErr(errCode.TokenInvalid)
			resp.Abort()
			return
		}

		if str != tokenStr {
			resp.ResponseErr(errCode.TokenInvalid)
			resp.Abort()
			return
		}

		handler(resp, tokenInfo)
	}
}
