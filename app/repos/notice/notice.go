package notice_repo

import (
	"vote/vote/app/db"
	"vote/vote/app/models"
	"vote/vote/app/models/api_model"
)

//创建公告
func Create(notice *models.Notice) error {
	err := db.MyDB.Model(&models.Notice{}).Create(&notice).Error
	if err != nil {
		return err
	}
	return nil
}

//获取一条公告信息
func Get(id string) (*models.Notice, error) {
	var notice models.Notice
	err := db.MyDB.Model(&models.Notice{}).Where("id=?", id).First(&notice).Error
	if err != nil {
		return nil, err
	}
	return &notice, nil
}

//修改公告信息
func Update(notice models.Notice) error {
	err := db.MyDB.Model(&models.Notice{}).Where("id=?", notice.ID).Update(&notice).Error
	if err != nil {
		return err
	}
	return nil
}

//删除一条公告信息
func Delete(id string) error {
	err := db.MyDB.Model(&models.Notice{}).Where("id=?", id).Delete(&models.Notice{}).Error
	if err != nil {
		return err
	}
	return nil
}

//获取公告列表
func ListForAdmin(limit, page int) (*api_model.AdminNoticeList, error) {
	var list api_model.AdminNoticeList

	query := db.MyDB.Model(models.Notice{}).
		Select("notice.id as id," +
			"notice.title as title," +
			"notice.abstract as abst," +
			"notice.created_time as created_time," +
			"notice.updated_time as updated_time," +
			"notice.initiator as author," +
			"notice.is_top as top").Count(&list.Total).Order("notice.is_top desc").Order("notice.created_time desc")

	err := query.Limit(limit).Offset(limit * (page - 1)).Scan(&list.List).Error

	if err != nil {
		return nil, err
	}

	return &list, nil
}

//用户端公告列表
func List(page, limit int) ([]api_model.NoticeList, error) {
	var list = make([]api_model.NoticeList, 0)

	if err := db.MyDB.Model(&models.Notice{}).
		Select("id as id," +
			"title as title," +
			"abstract as abst," +
			"created_time as time," +
			"initiator  as author," +
			"is_top as top").
		Order("is_top desc").
		Order("created_time desc").
		Limit(limit).Offset(limit * (page - 1)).
		Scan(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}
