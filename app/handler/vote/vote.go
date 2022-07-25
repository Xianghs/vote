package vote_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"vote/vote/app/errCode"
	"vote/vote/app/models"
	vote_service "vote/vote/app/services/vote"
	"vote/vote/app/util/admin_token"
	"vote/vote/app/util/resp"
	"vote/vote/app/util/token"
)

//创建投票
func Create(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var (
		form struct {
			Title   string   `json:"title" binding:"required"`
			Author  string   `json:"author" binding:"required"`
			Content string   `json:"content" binding:"required"`
			EndTime string   `json:"end_time" binding:"required"`
			Img     string   `json:"img" `
			Items   []string `json:"items" binding:"required"`
		}
		vote models.Vote
	)

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	vote.Title = form.Title
	vote.Initiator = form.Author
	vote.Description = form.Content
	end_t, _ := time.ParseInLocation("2006-01-02 15:04:05", form.EndTime, time.Local)
	vote.EndTime = end_t
	vote.Img = form.Img

	err := vote_service.Create(form.Items, vote)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()
}

//统计
func Statistics(resp *resp.Resp, tokenInfo *admin_token.Claims) {

	count, err := vote_service.Statistics()
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.Response(count)
}

//管理员页面投票列表
func List(c *gin.Context) {
	var (
		form struct {
			Key   string `json:"key"`
			Tag   string `json:"tag"`
			Limit int    `json:"limit" binding:"required"`
			Page  int    `json:"page" binding:"required"`
		}
		resp1 = resp.NewResp(c)
	)

	if err := resp1.ShouldBind(&form); err != nil {
		resp1.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	list, err := vote_service.ListForAdmin(form.Key, form.Tag, form.Page, form.Limit)

	if err != nil {
		resp1.ResponseErr(errCode.OtherError(err))
		return
	}

	resp1.Response(list)
}

//删除投票信息
func Delete(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var form struct {
		ID string `json:"id" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	err := vote_service.Delete(form.ID)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()
}

//投票主题详情
func Detail(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var form struct {
		ID string `json:"id" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	detail, err := vote_service.Detail(form.ID)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.Response(detail)
}

//修改投票主题
func Update(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var (
		form struct {
			ID      string `json:"id" binding:"required"`
			Title   string `json:"title" binding:"required"`
			Author  string `json:"author" binding:"required"`
			Content string `json:"content" binding:"required"`
			EndTime string `json:"end_time" binding:"required"`
			State   string `json:"state" binding:"required"`
			Img     string `json:"img" `
		}
		vote models.Vote
	)

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	end_t, _ := time.ParseInLocation("2006-01-02 15:04:05", form.EndTime, time.Local)

	vote.ID = form.ID
	vote.Title = form.Title
	vote.Initiator = form.Author
	vote.Description = form.Content
	vote.EndTime = end_t
	vote.Status = form.State
	vote.Img = form.Img
	vote.UpdatedTime = time.Now()

	err := vote_service.Update(vote)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()
}

//投票结果
func Result(c *gin.Context) {
	var form struct {
		ID string `json:"id" binding:"required"`
	}

	var resp1 = resp.NewResp(c)

	if err := resp1.ShouldBind(&form); err != nil {
		resp1.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	options, err := vote_service.Result(form.ID)

	if err != nil {
		resp1.ResponseErr(errCode.OtherError(err))
		return
	}

	resp1.Response(options)
}

//用户端投票详情，附带我的选择
func DetailForUser(resp *resp.Resp, token *token.Claims) {
	var form struct {
		ID string `json:"id" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	detail, err := vote_service.DetailForUser(form.ID, token.ID)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.Response(detail)
}

//我的投票列表
func MyVotes(resp *resp.Resp, token *token.Claims) {

	list, err := vote_service.MyVotes(token.ID)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.Response(list)
}

//投票操作
func Vote(resp *resp.Resp, token *token.Claims) {
	var (
		form struct {
			ID     string `json:"id" binding:"required"`
			Choice int    `json:"choice" `
		}
		record models.VoteRecord
	)

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		fmt.Println(err.Error())
		return
	}

	record.UserID = token.ID
	record.VoteID = form.ID
	record.Choice = form.Choice
	record.VoteTime = time.Now()

	if err := vote_service.Vote(record); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()

}
