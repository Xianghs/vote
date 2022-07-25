package resp

import (
	"github.com/gin-gonic/gin"
	"vote/vote/app/errCode"
)

type Resp struct {
	*gin.Context
}

func NewResp(c *gin.Context) *Resp {
	return &Resp{c}
}

func (resp *Resp) ResponseErr(errCode errCode.ErrorCode) {
	jsonMap := map[string]interface{}{
		"Code":    errCode.Code,
		"Message": errCode.Message,
		"data":    "",
	}
	resp.send(200, jsonMap)
}

func (resp *Resp) ResponseOK() {
	jsonMap := map[string]interface{}{
		"Code":    200,
		"Message": "OK",
		"Data":    "",
	}
	resp.send(200, jsonMap)
}

func (resp *Resp) Response(data interface{}) {
	jsonMap := map[string]interface{}{
		"Code":    200,
		"Message": "OK",
		"Data":    data,
	}
	resp.send(200, jsonMap)
}

//发送
func (resp *Resp) send(code int, m map[string]interface{}) {
	resp.JSON(code, m)
}
