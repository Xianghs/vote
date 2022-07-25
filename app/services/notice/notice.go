package notice_service

import (
	"time"
	"vote/vote/app/models"
	"vote/vote/app/models/api_model"
	notice_repo "vote/vote/app/repos/notice"
	"vote/vote/app/util"
)

//发布公告
func Post(notice *models.Notice) error {
	notice.CreatedTime = time.Now()
	notice.UpdatedTime = time.Now()
	notice.ID = util.CreateRandomString(20)

	//提取摘要
	str := notice.Description
	if len(str) > 90 {
		notice.Abstract = str[0:90] + "..."
	} else {
		notice.Abstract = str
	}

	//创建公告记录
	err := notice_repo.Create(notice)
	if err != nil {
		return err
	}

	return nil
}

//后台公告列表
func ListForAdmin(limit, page int) (*api_model.AdminNoticeList, error) {
	list, err := notice_repo.ListForAdmin(limit, page)

	if err != nil {
		return nil, err
	}

	return list, nil
}

//删除一条公告
func Delete(id string) error {
	err := notice_repo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

//公告详情
func Detail(id string) (*models.Notice, error) {
	notice, err := notice_repo.Get(id)

	if err != nil {
		return nil, err
	}

	return notice, nil
}

//编辑公告
func Update(notice models.Notice) error {
	n, err := notice_repo.Get(notice.ID)
	if err != nil {
		return err
	}

	n.Title = notice.Title
	n.Initiator = notice.Initiator
	n.Description = notice.Description
	if notice.Img != "" {
		n.Img = notice.Img
	}
	n.Abstract = notice.Abstract
	n.IsTop = notice.IsTop

	err = notice_repo.Update(*n)

	if err != nil {
		return err
	}

	return nil
}

//用户端公告列表
func List(page, limit int) ([]api_model.NoticeList, error) {
	list, err := notice_repo.List(page, limit)

	if err != nil {
		return nil, err
	}

	return list, nil
}
