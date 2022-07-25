package admin_service

import (
	"errors"
	"vote/vote/app/models"
	admin_repo "vote/vote/app/repos/admin"
	"vote/vote/app/util"
	"vote/vote/app/util/admin_token"
	"vote/vote/app/validators"
)

//管理员登陆
func SignIn(adm models.Admin) (string, error) {
	//验证特殊字符
	if validators.IsSpecialChar(adm.AdminName) {
		return "", errors.New("名称中不能含有特殊字符！")
	}

	//获取该用户信息
	admin, err1 := admin_repo.Get(adm.AdminName)
	if err1 != nil {
		return "", err1
	}
	pwd, e := util.AesEncrypt(adm.Password)
	if e != nil {
		return "", e
	}

	if pwd != admin.Password {
		return "", errors.New("登录失败，密码错误")
	}

	auth := admin_token.GenerateToken(*admin)

	return auth, nil
}

//修改密码
func Reset(prepwd, newpwd string, token *admin_token.Claims) error {

	var admin *models.Admin
	var err error

	admin, err = admin_repo.GetForID(token.ID)
	if err != nil {
		return err
	}

	pre_pwd, _ := util.AesEncrypt(prepwd)
	new_pwd, _ := util.AesEncrypt(newpwd)

	if pre_pwd != admin.Password {
		return errors.New("原密码错误！")
	}

	admin.Password = new_pwd

	//更新
	err1 := admin_repo.Update(*admin)
	if err1 != nil {
		return err1
	}
	return nil
}

//获取自己的信息
func Info(token *admin_token.Claims) (*models.Admin, error) {
	admin, err := admin_repo.GetForID(token.ID)
	if err != nil {
		return nil, err
	}
	return admin, err
}

//获取系统信息
func SystemInfo() (*models.SystemSet, error) {
	sys, err := admin_repo.SystemInfo()

	if err != nil {
		return nil, err
	}

	return sys, nil
}

//修改系统信息
func UpdateSystem(set models.SystemSet) error {

	set.ID = 1

	err := admin_repo.UpdateSystem(set)

	if err != nil {
		return err
	}

	return nil
}
