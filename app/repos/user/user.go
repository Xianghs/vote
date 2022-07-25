package user_repo

import (
	"vote/vote/app/db"
	"vote/vote/app/models"
	"vote/vote/app/models/api_model"
)

//创建用户
func Create(user models.User) error {
	err := db.MyDB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

//获取一个用户
func Get(phone string) (*models.User, error) {
	var user models.User
	err := db.MyDB.Model(&models.User{}).Where("phone=?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetByID(ID string) (*models.User, error) {
	var user models.User
	err := db.MyDB.Model(&models.User{}).Where("id=?", ID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//修改
func Update(user models.User) error {
	err := db.MyDB.Model(&models.User{}).Where("id = ?", user.ID).Update(&user).Error
	if err != nil {
		return err
	}
	return nil
}

//获取用户列表
func List(limit, page int, tag uint) (*api_model.UserList, error) {
	var list api_model.UserList

	query := db.MyDB.Model(&models.User{}).
		Select("user.id as id,user.name as name,user.head_img as head,user.phone as phone,user.signup_time as createdtime,user.auth as auth")

	if tag != 0 && tag <= 4 {
		query = query.Where("auth=?", tag)
	}

	query.Count(&list.Total)

	err := query.Order("signup_time desc").Limit(limit).Offset(limit * (page - 1)).Scan(&list.List).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

//删除一个用户
func Delete(id string) error {
	err := db.MyDB.Model(&models.User{}).Where("id=?", id).Delete(&models.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

//搜索
func Search(key string) (*api_model.UserList, error) {
	var list api_model.UserList

	if err := db.MyDB.Model(&models.User{}).
		Select("user.id as id,user.name as name,user.head_img as head,user.phone as phone,user.signup_time as createdtime,user.auth as auth").
		Where("name like ?", "%"+key+"%").Scan(&list.List).Count(&list.Total).Error; err != nil {
		return nil, err
	}

	return &list, nil
}
