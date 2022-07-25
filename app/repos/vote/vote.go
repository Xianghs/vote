package vote_repo

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"vote/vote/app/db"
	"vote/vote/app/models"
	"vote/vote/app/models/api_model"
)

//创建
func Create(vote api_model.VoteDetail) error {
	tx := db.MyDB.Begin()

	if err := tx.Model(&models.Vote{}).Create(&vote.Info).Error; err != nil {
		tx.Rollback()
		return err
	}

	//拼接SQL，批量插入选项记录
	sql := "INSERT INTO `vote_item` (`vote_id`,`item_id`,`item_name`) VALUES "
	for k, v := range vote.Item {
		if len(vote.Item)-1 == k {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("('%s',%d,'%s');", v.VoteID, v.ItemID, v.ItemName)
		} else {
			sql += fmt.Sprintf("('%s',%d,'%s'),", v.VoteID, v.ItemID, v.ItemName)
		}
	}

	if err := tx.Exec(sql).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//删除
func Delete(id string) error {
	tx := db.MyDB.Begin()

	if err := tx.Model(&models.Vote{}).Where("id=?", id).Delete(&models.Vote{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.VoteItem{}).Where("vote_id=?", id).Delete(&models.VoteItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.VoteRecord{}).Where("vote_id=?", id).Delete(&models.VoteRecord{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//修改
func Update(vote models.Vote) error {
	if err := db.MyDB.Model(&models.Vote{}).Where("id=?", vote.ID).Update(&vote).Error; err != nil {
		return err
	}

	return nil
}

//获取一个投票信息
func Get(id string) (*models.Vote, error) {
	var vote models.Vote

	if err := db.MyDB.Model(&models.Vote{}).Where("id=?", id).First(&vote).Error; err != nil {
		return nil, err
	}

	return &vote, nil
}

//获取一个投票的所有选项
func GetItem(id string) ([]models.VoteItem, error) {
	var items = make([]models.VoteItem, 0)

	if err := db.MyDB.Model(&models.VoteItem{}).Where("vote_id=?", id).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

//统计
func Statistics() (*api_model.VoteStatistics, error) {
	var count api_model.VoteStatistics

	if err := db.MyDB.Model(&models.Vote{}).Count(&count.AllCount).Error; err != nil {
		return nil, err
	}

	if err := db.MyDB.Model(&models.Vote{}).Where("status='进行中'").Count(&count.RunningCount).Error; err != nil {
		return nil, err
	}

	if err := db.MyDB.Model(&models.Vote{}).Where("status='已结束'").Count(&count.OverCount).Error; err != nil {
		return nil, err
	}

	if err := db.MyDB.Model(&models.VoteRecord{}).Count(&count.VotedCount).Error; err != nil {
		return nil, err
	}

	return &count, nil
}

//管理员页面投票列表
func ListForAdmin(key, tag string, page, limit int) (*api_model.AdminVoteList, error) {
	var list api_model.AdminVoteList

	query := db.MyDB

	if tag != "" {
		query = query.Where("status=?", tag)
	}

	if key != "" {
		query = query.Where("vote.title like ?", "%"+key+"%")
	}

	if err := query.Model(&models.Vote{}).Count(&list.Total).
		Select("count(vote_record.choice) as peoples," +
			"vote.id as id," +
			"vote.title as title," +
			"vote.status as state," +
			"vote.created_time as starttime," +
			"vote.end_time as endtime").
		Joins("left join vote_record on vote_record.vote_id=vote.id").
		Group("vote.id").
		Order("vote.created_time desc").
		Limit(limit).Offset(limit * (page - 1)).
		Scan(&list.List).Error; err != nil {
		return nil, err
	}

	return &list, nil
}

//投票结果
func Result(id string) ([]api_model.VoteOptions, error) {
	var (
		result   = make([]api_model.VoteOptions, 0)
		options  = make([]models.VoteItem, 0)
		system   models.SystemSet
		allCount = 0
		err      error
	)

	//得到总户数和总面积数
	if err = db.MyDB.Model(&models.SystemSet{}).Where("id=1").First(&system).Error; err != nil {
		return nil, err
	}

	//得到所有选项
	if err = db.MyDB.Model(&models.VoteItem{}).
		Where("vote_id=?", id).
		Order("item_id asc").
		Find(&options).Error; err != nil {
		return nil, err
	}

	//得到所有选项总票数
	if err = db.MyDB.Model(&models.VoteRecord{}).Where("vote_id=?", id).Count(&allCount).Error; err != nil {
		return nil, err
	}

	for _, v := range options {
		count := 0
		//获取该选项票数
		if err = db.MyDB.Model(&models.VoteRecord{}).
			Where("vote_id=? and choice=?", id, v.ItemID).
			Count(&count).Error; err != nil {
			return nil, err
		}
		//获取该选项的面积总和
		areaSum := struct {
			SumArea float64 `json:"sum_area"`
		}{SumArea: 0}
		if err = db.MyDB.Model(&models.VoteRecord{}).
			Select("sum(user.occupied_area) as sum_area").
			Joins("left join user on user.id=vote_record.user_id").
			Where("vote_record.vote_id=? and vote_record.choice=?", id, v.ItemID).
			Scan(&areaSum).Error; err != nil {
			return nil, err
		}

		proportion := int((float64(count) / float64(allCount)) * 100)

		if proportion < 0 {
			proportion = 0
		}

		//拼装返回信息
		result = append(result, api_model.VoteOptions{
			ID:         v.ItemID,
			Name:       v.ItemName,
			Count:      count,
			Area:       areaSum.SumArea,
			Proportion: proportion,
			Areap:      int((areaSum.SumArea / system.TotalArea) * 100),
		})
	}

	return result, nil
}

//用户端投票详情，附带我的选择
func DetailForUser(voteId, userId string) (*api_model.VoteDetailForUser, error) {
	var detail api_model.VoteDetailForUser

	vote, err := Get(voteId)

	if err != nil {
		return nil, err
	}

	detail.ID = vote.ID
	detail.Img = vote.Img
	detail.State = vote.Status
	detail.Content = vote.Description
	detail.Author = vote.Initiator
	detail.Title = vote.Title
	detail.EndT = vote.EndTime
	detail.StartT = vote.CreatedTime
	detail.UpdatedAt = vote.UpdatedTime

	count := 0
	choose := struct {
		Choice int `json:"choice"`
	}{Choice: 0}

	if err := db.MyDB.Model(&models.VoteRecord{}).Select("choice as choice").
		Where("vote_id=? and user_id=?", voteId, userId).
		Count(&count).
		Scan(&choose).
		Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if count == 0 {
		detail.MyChoose = -1
		return &detail, nil
	}

	detail.MyChoose = choose.Choice

	return &detail, nil
}

//我的投票列表
func MyVotes(userId string) ([]api_model.MyVoteList, error) {
	var list = make([]api_model.MyVoteList, 0)

	if err := db.MyDB.Model(&models.Vote{}).
		Select("vote.id as id,"+
			"vote.title as title,"+
			"vote.status as state,"+
			"vote.created_time as starttime,"+
			"vote.end_time as endtime,"+
			"vote_record.choice as my_choose").
		Joins("left join vote_record on vote_record.vote_id=vote.id").
		Where("vote_record.user_id=?", userId).
		Order("vote_record.vote_time desc").
		Scan(&list).Error; err != nil {
		return nil, err
	}

	for k, v := range list {
		items := make([]models.VoteItem, 0)
		if err := db.MyDB.Model(&models.VoteItem{}).Where("vote_id=?", v.ID).Find(&items).Error; err != nil {
			return nil, err
		}

		for _, v1 := range items {
			choose, _ := strconv.Atoi(v.MyChoose)
			if choose == v1.ItemID {
				list[k].MyChoose = v1.ItemName
			}
		}
	}

	return list, nil

}

//投票操作
func Vote(voteRecord models.VoteRecord) error {

	//先判断是否已经在该主题投票
	count := 0
	if err := db.MyDB.Model(&models.VoteRecord{}).
		Where("vote_id=? and user_id=?", voteRecord.VoteID, voteRecord.UserID).
		Count(&count).Error; err != nil {
		return nil
	}

	//没有投过票直接创建，投过票的更新
	if count == 0 {
		if err := db.MyDB.Model(&models.VoteRecord{}).Create(&voteRecord).Error; err != nil {
			return err
		}
	} else {
		if err := db.MyDB.Model(&models.VoteRecord{}).
			Where("vote_id=? and user_id=?", voteRecord.VoteID, voteRecord.UserID).
			Updates(map[string]interface{}{"vote_id": voteRecord.VoteID, "user_id": voteRecord.UserID, "choice": voteRecord.Choice, "vote_time": voteRecord.VoteTime}).
			Error; err != nil {
			return err
		}
	}

	return nil

}
