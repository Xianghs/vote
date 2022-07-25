package vote_service

import (
	"errors"
	"time"
	"vote/vote/app/models"
	"vote/vote/app/models/api_model"
	user_repo "vote/vote/app/repos/user"
	vote_repo "vote/vote/app/repos/vote"
	"vote/vote/app/util"
)

func Detail(id string) (*api_model.VoteDetail, error) {

	var detail api_model.VoteDetail

	vote, err := vote_repo.Get(id)
	if err != nil {
		return nil, err
	}
	detail.Info = *vote

	items, err := vote_repo.GetItem(id)
	if err != nil {
		return nil, err
	}
	detail.Item = items

	return &detail, nil
}

//创建
func Create(itemName []string, info models.Vote) error {
	var vote api_model.VoteDetail

	vote.Item = make([]models.VoteItem, 0)

	vote.Info = info
	vote.Info.ID = util.CreateRandomString(20)
	vote.Info.Status = "进行中"
	vote.Info.CreatedTime = time.Now()
	vote.Info.UpdatedTime = time.Now()

	for k, v := range itemName {
		vote.Item = append(vote.Item, models.VoteItem{
			VoteID:   vote.Info.ID,
			ItemID:   k,
			ItemName: v,
		})
	}

	err := vote_repo.Create(vote)

	if err != nil {
		return err
	}

	return nil
}

//统计
func Statistics() (*api_model.VoteStatistics, error) {
	count, err := vote_repo.Statistics()

	if err != nil {
		return nil, err
	}

	return count, nil

}

//管理员页面投票列表,包含搜索
func ListForAdmin(key, tag string, page, limit int) (*api_model.AdminVoteList, error) {
	list, err := vote_repo.ListForAdmin(key, tag, page, limit)

	if err != nil {
		return nil, err
	}

	return list, nil
}

//删除投票信息
func Delete(id string) error {
	if err := vote_repo.Delete(id); err != nil {
		return err
	}

	return nil
}

func Update(vote models.Vote) error {
	if err := vote_repo.Update(vote); err != nil {
		return err
	}

	return nil
}

//投票结果
func Result(id string) ([]api_model.VoteOptions, error) {
	options, err := vote_repo.Result(id)

	if err != nil {
		return nil, err
	}

	return options, nil
}

//用户端投票详情，附带我的选择
func DetailForUser(voteId, userId string) (*api_model.VoteDetailForUser, error) {
	detail, err := vote_repo.DetailForUser(voteId, userId)

	if err != nil {
		return nil, err
	}

	return detail, nil
}

//我的投票列表
func MyVotes(userId string) ([]api_model.MyVoteList, error) {
	list, err := vote_repo.MyVotes(userId)

	if err != nil {
		return nil, err
	}

	return list, nil
}

//投票操作
func Vote(record models.VoteRecord) error {
	//先判断用户权限
	user, err := user_repo.GetByID(record.UserID)
	if err != nil {
		return err
	}

	if user.Auth != 3 {
		return errors.New("请先认证后再来投票！")
	}

	if err := vote_repo.Vote(record); err != nil {
		return err
	}

	return nil
}
