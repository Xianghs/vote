package admin_repo

import (
	"errors"
	"vote/vote/app/db"
	"vote/vote/app/models"
)

//获取一个管理员
func Get(name string) (*models.Admin, error) {
	var admin models.Admin
	err := db.MyDB.Model(&models.Admin{}).Where("admin_name=?", name).First(&admin).Error
	if err != nil {
		return nil, err
	}
	if admin.ID == "" {
		return nil, errors.New("没有该用户")
	}
	return &admin, nil
}
func GetForID(id string) (*models.Admin, error) {
	var admin models.Admin
	err := db.MyDB.Model(&models.Admin{}).Where("id=?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

//修改管理员信息
func Update(admin models.Admin) error {
	err := db.MyDB.Model(&models.Admin{}).Where("id=?", admin.ID).Update(&admin).Error
	if err != nil {
		return err
	}
	return nil
}

//获取系统信息
func SystemInfo() (*models.SystemSet, error) {
	var sys models.SystemSet

	if err := db.MyDB.Model(&models.SystemSet{}).Where("id=1").First(&sys).Error; err != nil {
		return nil, err
	}

	return &sys, nil
}

//修改系统信息
func UpdateSystem(sys models.SystemSet) error {
	if err := db.MyDB.Model(&models.SystemSet{}).Where("id=1").Update(&sys).Error; err != nil {
		return err
	}

	return nil
}
