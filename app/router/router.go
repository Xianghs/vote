package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"vote/vote/app/config"
	admin_handler "vote/vote/app/handler/admin"
	notice_handler "vote/vote/app/handler/notice"
	user_handler "vote/vote/app/handler/user"
	vote_handler "vote/vote/app/handler/vote"
	"vote/vote/app/middleware"
	"vote/vote/app/util/admin_token"
	"vote/vote/app/util/token"
)

func InitRouter(router *gin.Engine) {
	//跨域
	router.Use(cors.New(config.Cors))
	//访问频率限制
	router.Use(middleware.LimitMiddleware())
	//日志
	router.Use(middleware.Logger())
	api(router)
}

func api(router *gin.Engine) {
	//用户
	var user = router.Group("/user")
	{
		user.POST("/sign_up", user_handler.SignUp)                                    //注册
		user.POST("/sign_in", user_handler.SignIn)                                    //登陆
		user.POST("/reset_pwd", token.RequiredToken(user_handler.ResetPwd))           //修改密码
		user.POST("/reset_name", token.RequiredToken(user_handler.UpdateName))        //修改用户名
		user.POST("/reset_phone", token.RequiredToken(user_handler.UpdatePhone))      //修改电话号码
		user.POST("/reset_head", token.RequiredToken(user_handler.UpdateHead))        //修改用户头像
		user.POST("/submit_auth", token.RequiredToken(user_handler.SubmitAuth))       //提交认证信息
		user.POST("/info", token.RequiredToken(user_handler.Info))                    //获取自己的信息
		user.POST("/list", admin_token.RequiredToken(user_handler.List))              //获取用户列表
		user.POST("/delete", admin_token.RequiredToken(user_handler.Delete))          //删除一个用户
		user.POST("/upload", token.RequiredToken(user_handler.Upload))                //用户上传图片
		user.POST("/search", admin_token.RequiredToken(user_handler.Search))          //搜索
		user.POST("/detail", admin_token.RequiredToken(user_handler.Detail))          //用户详情
		user.POST("/update", admin_token.RequiredToken(user_handler.Update))          //管理员直接编辑用户信息
		user.POST("/update_auth", admin_token.RequiredToken(user_handler.UpdateAuth)) //管理员修改用户的权限级别
	}

	//管理员
	var admin = router.Group("/admin")
	{
		admin.POST("/sign_in", admin_handler.SignIn)                                        //管理员登陆
		admin.POST("/reset", admin_token.RequiredToken(admin_handler.Reset))                //修改管理员密码
		admin.POST("/upload", admin_token.RequiredToken(admin_handler.Upload))              //管理员上传文件
		admin.POST("/system_info", admin_token.RequiredToken(admin_handler.SystemInfo))     //获取系统信息
		admin.POST("/system_update", admin_token.RequiredToken(admin_handler.UpdateSystem)) //修改系统信息

	}

	//公告信息
	var notice = router.Group("/notice")
	{
		notice.POST("/post", admin_token.RequiredToken(notice_handler.Post))                   //发布公告
		notice.POST("/list_for_admin", admin_token.RequiredToken(notice_handler.ListForAdmin)) //后台公告列表
		notice.POST("/delete", admin_token.RequiredToken(notice_handler.Delete))               //删除一条公告
		notice.POST("/detail", notice_handler.Detail)                                          //公告详情
		notice.POST("/update", admin_token.RequiredToken(notice_handler.Update))               //修改公告
		notice.POST("/list", notice_handler.List)                                              //用户端公告列表
	}

	//投票
	var vote = router.Group("/vote")
	{
		vote.POST("/create", admin_token.RequiredToken(vote_handler.Create))           //创建投票
		vote.POST("/count", admin_token.RequiredToken(vote_handler.Statistics))        //统计
		vote.POST("/list_for_admin", vote_handler.List)                                //投票列表,包含搜索功能
		vote.POST("/update", admin_token.RequiredToken(vote_handler.Update))           //修改投票信息
		vote.POST("/result", vote_handler.Result)                                      //投票结果
		vote.POST("/detail", admin_token.RequiredToken(vote_handler.Detail))           //管理员端投票主题详情
		vote.POST("/detail_for_user", token.RequiredToken(vote_handler.DetailForUser)) //用户端投票详情，附带自己的选择
		vote.POST("/my_votes", token.RequiredToken(vote_handler.MyVotes))              //我的投票列表
		vote.POST("/vote", token.RequiredToken(vote_handler.Vote))                     //投票操作
		vote.POST("/delete", admin_token.RequiredToken(vote_handler.Delete))           //删除投票主题
	}
}
