package notice_handler

import (
	"github.com/gin-gonic/gin"
	"vote/vote/app/errCode"
	"vote/vote/app/models"
	notice_service "vote/vote/app/services/notice"
	"vote/vote/app/util/admin_token"
	"vote/vote/app/util/resp"
)

//发布公告
func Post(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var (
		form struct {
			Title     string `json:"title" binding:"required"`
			Initiator string `json:"initiator" binding:"required"`
			Content   string `json:"description" binding:"required"`
			State     uint   `json:"state" binding:"required"`
			Img       string `json:"img" `
		}
		notice models.Notice
	)

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	notice.Title = form.Title
	notice.Description = form.Content
	notice.Initiator = form.Initiator
	notice.IsTop = int(form.State)
	notice.Img = form.Img

	err := notice_service.Post(&notice)
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}
	resp.ResponseOK()
}

//后台公告列表
func ListForAdmin(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var form struct {
		Limit int `json:"limit"  binding:"required"`
		Page  int `json:"page" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	list, err := notice_service.ListForAdmin(form.Limit, form.Page)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.Response(list)
}

//删除一条公告
func Delete(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var form struct {
		ID string `json:"id" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	err := notice_service.Delete(form.ID)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()
}

//公告详情
func Detail(c *gin.Context) {
	var (
		form struct {
			ID string `json:"id" binding:"required"`
		}

		resp1 = resp.NewResp(c)
	)

	if err := resp1.ShouldBind(&form); err != nil {
		resp1.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	notice, err := notice_service.Detail(form.ID)

	if err != nil {
		resp1.ResponseErr(errCode.OtherError(err))
		return
	}

	resp1.Response(notice)

}

//修改公告信息
func Update(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var (
		form struct {
			ID      string `json:"id" binding:"required"`
			Title   string `json:"title" binding:"required"`
			Author  string `json:"author" binding:"required"`
			Content string `json:"content" binding:"required"`
			Img     string `json:"img"`
			Abst    string `json:"abst" binding:"required"`
			Top     int    `json:"top" binding:"required"`
		}
		notice models.Notice
	)

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	notice.ID = form.ID
	notice.Title = form.Title
	notice.Initiator = form.Author
	notice.Description = form.Content
	notice.Img = form.Img
	notice.Abstract = form.Abst
	notice.IsTop = form.Top

	err := notice_service.Update(notice)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()
}

//用户端公告列表
func List(c *gin.Context) {
	var (
		form struct {
			Page  int `json:"page" binding:"required"`
			Limit int `json:"limit" binding:"required"`
		}

		resp1 = resp.NewResp(c)
	)

	if err := resp1.ShouldBind(&form); err != nil {
		resp1.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	list, err := notice_service.List(form.Page, form.Limit)

	if err != nil {
		resp1.ResponseErr(errCode.OtherError(err))
		return
	}

	resp1.Response(list)

}
