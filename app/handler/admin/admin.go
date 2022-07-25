package admin_handler

import (
	"github.com/gin-gonic/gin"
	"vote/vote/app/errCode"
	"vote/vote/app/models"
	admin_service "vote/vote/app/services/admin"
	"vote/vote/app/util"
	"vote/vote/app/util/admin_token"
	"vote/vote/app/util/resp"
)

//后台登陆
func SignIn(c *gin.Context) {
	var (
		form struct {
			Name     string `json:"name" binding:"required"`
			PassWord string `json:"pass_word" binding:"required"`
		}
		resp1 = resp.NewResp(c)
		admin models.Admin
	)
	if err := resp1.ShouldBind(&form); err != nil {
		resp1.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	admin.AdminName = form.Name
	admin.Password = form.PassWord
	auth, err := admin_service.SignIn(admin)
	if err != nil {
		resp1.ResponseErr(errCode.OtherError(err))
		return
	}
	resp1.Response(gin.H{"token": auth})
}

//修改信息
func Reset(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var form struct {
		PrePwd string `json:"pre_pwd" binding:"required"`
		NewPwd string `json:"new_pwd" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	err := admin_service.Reset(form.PrePwd, form.NewPwd, tokenInfo)
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()
}

//上传图片
func Upload(resp *resp.Resp, tokenInfo *admin_token.Claims) {

	file, err := resp.FormFile("file")

	if err != nil {
		resp.ResponseErr(errCode.FileUploadFail)
		return
	}

	path, err := util.UploadFile(file, tokenInfo.ID, "image/jpg")

	if err != nil {
		resp.ResponseErr(errCode.FileUploadFail)
		return
	}

	resp.Response(gin.H{"path": path})
}

//获取系统信息
func SystemInfo(resp *resp.Resp, tokenInfo *admin_token.Claims) {

	sys, err := admin_service.SystemInfo()

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.Response(sys)
}

//修改系统信息
func UpdateSystem(resp *resp.Resp, tokenInfo *admin_token.Claims) {
	var (
		form struct {
			TotalHousehold int     `json:"total_household" binding:"required"`
			TotalArea      float64 `json:"total_area" binding:"required"`
		}

		sys models.SystemSet
	)

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	sys.TotalHousehold = form.TotalHousehold
	sys.TotalArea = form.TotalArea

	err := admin_service.UpdateSystem(sys)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()
}
