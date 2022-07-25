package errCode

import "net/http"

type ErrorCode struct {
	Code     int    `Description:"自定义的错误码"`
	Message  string `Description:"错误码对应的消息"`
	HttpCode int    `Description:"标准 http Code"`
}

var (
	TooManyRequests = ErrorCode{
		Code:     http.StatusTooManyRequests,
		Message:  "您的请求过于频繁，请稍候再试",
		HttpCode: 429,
	}

	TokenInvalid = ErrorCode{
		Code:     700,
		Message:  "token失效,请重新登陆",
		HttpCode: 700,
	}

	AdminAuthInvalid = ErrorCode{
		Code:     800,
		Message:  "你不是管理员，不能进行该操作",
		HttpCode: 800,
	}

	FileUploadFail = ErrorCode{
		Code:     900,
		Message:  "文件上传失败",
		HttpCode: 900,
	}
)

func ErrorInvalidParameters(err error) ErrorCode {
	return ErrorCode{
		Code:     600,
		Message:  "参数错误：" + err.Error(),
		HttpCode: 600,
	}
}

func OtherError(err error) ErrorCode {
	return ErrorCode{
		Code:     10000,
		Message:  err.Error(),
		HttpCode: 500,
	}
}
