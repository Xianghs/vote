package validators

import (
	"regexp"
)

//自定义电话号码验证器
func PhoneValidator(phone string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)

}

//自定义身份证号码验证器
func IDCardNumValidator(num string) bool {
	regular := `^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(num)
}

//是否含有特殊字符
func IsSpecialChar(str string) bool {
	regular := "[ _`~!@#$%^&*()+=|{}':;',\\[\\].<>/?~！@#￥%……&*（）——+|{}【】‘；：”“’。，、？]|\n|\r|\t"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}
