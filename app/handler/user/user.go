package user_handler

import (
	"github.com/gin-gonic/gin"
	"vote/vote/app/errCode"
	"vote/vote/app/models"
	"vote/vote/app/util"
	"vote/vote/app/util/admin_token"
	"vote/vote/app/util/resp"
	"vote/vote/app/util/token"

	user_service "vote/vote/app/services/user"
)

//注册
func SignUp(c *gin.Context) {
	var (
		form struct {
			Phone    string `json:"phone" binding:"required"`
			PassWord string `json:"pass_word" binding:"required"`
		}
		resp1 = resp.NewResp(c)
		user  models.User
	)
	if err := resp1.ShouldBind(&form); err != nil {
		resp1.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	user.Phone = form.Phone
	user.PassWord = form.PassWord
	err := user_service.SignUp(user)
	if err != nil {
		resp1.ResponseErr(errCode.OtherError(err))
		return
	}
	resp1.ResponseOK()
}

//登陆
func SignIn(c *gin.Context) {
	var (
		form struct {
			Phone    string `json:"phone" binding:"required"`
			PassWord string `json:"pass_word" binding:"required"`
		}
		resp1 = resp.NewResp(c)
		user  models.User
	)
	if err := resp1.ShouldBind(&form); err != nil {
		resp1.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}

	user.Phone = form.Phone
	user.PassWord = form.PassWord
	tokenInfo, err := user_service.SignIn(user)
	if err != nil {
		resp1.ResponseErr(errCode.OtherError(err))
		return
	}
	resp1.Response(gin.H{"token": tokenInfo})
}

//修改密码
func ResetPwd(resp *resp.Resp, token *token.Claims) {
	var form struct {
		PrePwd string `json:"pre_pwd" binding:"required"`
		NewPwd string `json:"new_pwd" binding:"required"`
	}
	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.ErrorInvalidParameters(err))
		return
	}
	err := user_service.ResetPwd(form.PrePwd, form.NewPwd, token.Phone)
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}
	resp.ResponseOK()
}

//修改用户名
func UpdateName(resp *resp.Resp, token *token.Claims) {
	var form struct {
		NewName string `json:"new_name" binding:"required"`
	}
	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}
	err := user_service.UpdateName(form.NewName, token)
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}
	resp.ResponseOK()
}

func UpdatePhone(resp *resp.Resp, token *token.Claims) {
	var form struct {
		NewPhone string `json:"new_name" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	err := user_service.UpdatePhone(form.NewPhone, token)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}
	resp.ResponseOK()
}

//提交认证信息
func SubmitAuth(resp *resp.Resp, token *token.Claims) {
	var (
		form struct {
			IDCard string  `json:"id_card" binding:"required"`
			Area   float64 `json:"area" binding:"required"`
			Proof  string  `json:"proof" binding:"required"`
		}
		user models.User
	)
	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	user.OccupiedArea = form.Area
	user.IDCard = form.IDCard
	user.Proof = form.Proof
	err := user_service.SubmitAuth(user, token)
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}
	resp.ResponseOK()
}

//获取自己的信息
func Info(resp *resp.Resp, token *token.Claims) {
	user, err := user_service.Info(token)
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}
	resp.Response(user)
}

//获取用户列表
func List(resp *resp.Resp, token *admin_token.Claims) {
	var form struct {
		Limit int  `json:"limit" binding:"required"`
		Page  int  `json:"page" binding:"required"`
		Tag   uint `json:"tag" `
	}
	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	list, err := user_service.List(form.Limit, form.Page, form.Tag)
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}
	resp.Response(list)
}

//删除一个用户
func Delete(resp *resp.Resp, token *admin_token.Claims) {
	var form struct {
		ID string `json:"id" binding:"required"`
	}
	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	err := user_service.Delete(form.ID)
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}
	resp.ResponseOK()
}

func Search(resp *resp.Resp, token *admin_token.Claims) {
	var form struct {
		Key string `json:"key" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	list, err := user_service.Search(form.Key)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.Response(list)
}

func Detail(resp *resp.Resp, token *admin_token.Claims) {
	var form struct {
		ID string `json:"id" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	user, err := user_service.Detail(form.ID)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.Response(user)
}

//管理员编辑用户信息
func Update(resp *resp.Resp, token *admin_token.Claims) {
	var (
		form struct {
			ID       string `json:"id" binding:"required"`
			Phone    string `json:"phone" `
			Name     string `json:"name" `
			Head     string `json:"head" `
			PassWord string `json:"pass_word" `
		}
		user models.User
	)

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	user.ID = form.ID
	user.Name = form.Name
	user.Phone = form.Phone
	user.PassWord = form.PassWord
	user.HeadImg = form.Head

	err := user_service.Update(user)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()
}

//修改用户权限
func UpdateAuth(resp *resp.Resp, token *admin_token.Claims) {
	var form struct {
		ID   string `json:"id" binding:"required"`
		Auth uint   `json:"auth" binding:"required"`
	}

	if err := resp.ShouldBind(&form); err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	err := user_service.UpdateAuth(form.ID, form.Auth)

	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.ResponseOK()
}

//修改头像
func UpdateHead(resp *resp.Resp, token *token.Claims) {
	file, err := resp.FormFile("file")

	if err != nil {
		resp.ResponseErr(errCode.FileUploadFail)
		return
	}

	path, err := util.UploadFile(file, token.ID, "image/jpg")

	if err != nil {
		resp.ResponseErr(errCode.FileUploadFail)
		return
	}

	err = user_service.UpdateHead(*path, token)
	if err != nil {
		resp.ResponseErr(errCode.OtherError(err))
		return
	}

	resp.Response(gin.H{"path": path})
}

//业主上传文件
func Upload(resp *resp.Resp, token *token.Claims) {
	file, err := resp.FormFile("file")

	if err != nil {
		resp.ResponseErr(errCode.FileUploadFail)
		return
	}

	path, err := util.UploadFile(file, token.ID, "image/jpg")

	if err != nil {
		resp.ResponseErr(errCode.FileUploadFail)
		return
	}

	resp.Response(gin.H{"path": path})
}
