package api_model

import (
	"time"
	"vote/vote/app/models"
)

//管理员界面用户列表
type UserList struct {
	List []struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Head        string    `json:"head"`
		Phone       string    `json:"phone"`
		Createdtime time.Time `json:"createdtime"`
		Auth        int       `json:"auth"`
	}

	Total int `json:"total"`
}

//管理员界面公告列表
type AdminNoticeList struct {
	List []struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Abst        string    `json:"abst"`
		CreatedTime time.Time `json:"createdtime"`
		UpdatedTime time.Time `json:"updatedtime"`
		Author      string    `json:"author"`
		Top         int       `json:"top"`
	}

	Total int `json:"total"`
}

//投票详情
type VoteDetail struct {
	Info models.Vote
	Item []models.VoteItem
}

//管理员页面投票相关统计
type VoteStatistics struct {
	AllCount     int `json:"all_count"`
	RunningCount int `json:"running_count"`
	OverCount    int `json:"over_count"`
	VotedCount   int `json:"voted_count"`
}

//管理员界面投票列表
type AdminVoteList struct {
	List []struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		State     string    `json:"state"`
		Peoples   int       `json:"peoples"` //参与人数
		Starttime time.Time `json:"starttime"`
		Endtime   time.Time `json:"endtime"`
	}
	Total int `json:"total"`
}

//投票选项及其得票统计
type VoteOptions struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Count      int     `json:"count"`      //票数
	Area       float64 `json:"area"`       //面积
	Proportion int     `json:"proportion"` //票数占比
	Areap      int     `json:"areap"`      //面积占比
}

//用户端主页公告列表
type NoticeList struct {
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	Abst   string    `json:"abst"`
	Time   time.Time `json:"time"`
	Author string    `json:"author"`
	Top    int       `json:"top"`
}

//用户端投票详情，附带我的选择
type VoteDetailForUser struct {
	ID        string    `json:"id"`
	State     string    `json:"state"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	Img       string    `json:"img"`
	MyChoose  int       `json:"my_choose"`
	StartT    time.Time `json:"start_t"`
	EndT      time.Time `json:"end_t"`
}

//我的投票列表
type MyVoteList struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`

	MyChoose  string    `json:"my_choose"`
	Starttime time.Time `json:"starttime"`
	Endtime   time.Time `json:"endtime"`
}
