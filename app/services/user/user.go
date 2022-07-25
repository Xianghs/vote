package user_service

import (
	"errors"
	"time"
	"vote/vote/app/db"
	"vote/vote/app/models"
	"vote/vote/app/models/api_model"
	user_repo "vote/vote/app/repos/user"
	"vote/vote/app/util/token"

	"vote/vote/app/util"
	"vote/vote/app/validators"
)

//注册
func SignUp(user models.User) error {
	//验证手机号格式
	if !validators.PhoneValidator(user.Phone) {
		return errors.New("请输入正确的电话号码！")
	}

	//验证密码中是否含有特殊字符
	if validators.IsSpecialChar(user.PassWord) {
		return errors.New("密码不能含有特殊字符！")
	}
	//判断电话号是否占用
	count := 0
	if err := db.MyDB.Model(&models.User{}).
		Where("phone=?", user.Phone).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return errors.New("电话号占用")
	}
	user.Name = "新业主"
	user.Auth = 1
	user.SignupTime = time.Now()
	user.SubmitTime = time.Date(1, 1, 1, 0, 0, 0, 0, time.Now().Location())
	user.AdoptTime = time.Date(1, 1, 1, 0, 0, 0, 0, time.Now().Location())
	user.ID = util.CreateRandomString(20)
	var e error
	user.PassWord, e = util.AesEncrypt(user.PassWord)
	if e != nil {
		return e
	}
	if err := user_repo.Create(user); err != nil {
		return err
	}
	return nil
}

//登陆
func SignIn(form models.User) (string, error) {
	//先判断手机号合法性
	if !validators.PhoneValidator(form.Phone) {
		return "", errors.New("手机号不合法，请重新输入")
	}
	//验证密码中是否含有特殊字符
	if validators.IsSpecialChar(form.PassWord) {
		return "", errors.New("密码不能含有特殊字符！")
	}
	//从数据库获取到该手机号用户
	user, err := user_repo.Get(form.Phone)
	if err != nil {
		return "", err
	}
	//解密
	password, err := util.AesDecrypt(user.PassWord)
	if err != nil {
		return "", errors.New("服务器繁忙,请稍后再试")
	}
	if password != form.PassWord {
		return "", errors.New("密码错误，登录失败")
	}

	//token生成
	auth := token.GenerateToken(*user)

	return auth, nil
}

//修改密码
func ResetPwd(pre_pwd, new_pwd, phone string) error {
	if validators.IsSpecialChar(pre_pwd) {
		return errors.New("密码不能含有特殊字符！")
	}
	if validators.IsSpecialChar(new_pwd) {
		return errors.New("密码不能含有特殊字符！")
	}
	user, err := user_repo.Get(phone)
	if user == nil || err != nil {
		return errors.New("服务器出错")
	}

	var e error
	pre_pwd, e = util.AesEncrypt(pre_pwd)
	if e != nil {
		return e
	}
	if pre_pwd != user.PassWord {
		return errors.New("原密码错误")
	}

	new_pwd, e = util.AesEncrypt(new_pwd)
	if e != nil {
		return e
	}
	user.PassWord = new_pwd
	err = user_repo.Update(*user)
	if err != nil {
		return err
	}
	return nil
}

//修改用户名
func UpdateName(newname string, token *token.Claims) error {
	if validators.IsSpecialChar(newname) {
		return errors.New("不能含有特殊字符！")
	}
	user, err := user_repo.Get(token.Phone)
	if err != nil || user == nil {
		return errors.New("服务器繁忙")
	}
	user.Name = newname
	err = user_repo.Update(*user)
	if err != nil {
		return err
	}
	return nil
}

//用户修改头像
func UpdateHead(newhead string, token *token.Claims) error {
	user, err := user_repo.Get(token.Phone)
	if err != nil {
		return errors.New("服务器繁忙，请稍后再试")
	}

	user.HeadImg = newhead

	err = user_repo.Update(*user)
	if err != nil {
		return err
	}

	return nil
}

//修改电话号码
func UpdatePhone(newphone string, token *token.Claims) error {
	//验证手机号格式
	if !validators.PhoneValidator(token.Phone) {
		return errors.New("请输入正确的电话号码！")
	}

	user, err := user_repo.Get(token.Phone)
	if err != nil || user == nil {
		return errors.New("服务器繁忙，请稍后再试！")
	}

	user.Phone = newphone
	err = user_repo.Update(*user)

	if err != nil {
		return err
	}

	return nil
}

//提交认证信息
func SubmitAuth(form models.User, token *token.Claims) error {

	if !validators.IDCardNumValidator(form.IDCard) {
		return errors.New("身份证号码格式错误")
	}
	user, err := user_repo.Get(token.Phone)
	if user == nil || err != nil {
		return errors.New("服务器繁忙")
	}

	user.IDCard = form.IDCard
	user.OccupiedArea = form.OccupiedArea
	user.Proof = form.Proof
	user.Auth = 2
	user.SubmitTime = time.Now()
	err = user_repo.Update(*user)
	if err != nil {
		return err
	}
	return nil
}

//获取自己的信息
func Info(token *token.Claims) (*models.User, error) {
	user, err := user_repo.Get(token.Phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//获取用户列表
func List(limit, page int, tag uint) (*api_model.UserList, error) {
	list, err := user_repo.List(limit, page, tag)
	if err != nil {
		return nil, err
	}
	return list, nil
}

//删除一个用户
func Delete(id string) error {
	err := user_repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

//搜索
func Search(key string) (*api_model.UserList, error) {
	list, err := user_repo.Search(key)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func Detail(ID string) (*models.User, error) {
	user, err := user_repo.GetByID(ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

//管理员编辑用户信息
func Update(newInfo models.User) error {
	user, err := user_repo.GetByID(newInfo.ID)

	if err != nil {
		return err
	}

	if newInfo.PassWord != "" {
		newpwd, err := util.AesEncrypt(newInfo.PassWord)
		user.PassWord = newpwd
		if err != nil {
			return err
		}
	}

	user.Name = newInfo.Name
	user.Phone = newInfo.Phone
	user.HeadImg = newInfo.HeadImg

	err = user_repo.Update(*user)
	if err != nil {
		return err
	}

	return nil
}

//修改用户权限
func UpdateAuth(ID string, auth uint) error {
	user, err := user_repo.GetByID(ID)

	if err != nil {
		return err
	}

	if auth > 4 {
		return errors.New("权限级别不合法")
	}

	if auth == 3 {
		user.AdoptTime = time.Now()
	}

	user.Auth = int(auth)

	err = user_repo.Update(*user)

	if err != nil {
		return err
	}

	return nil
}
