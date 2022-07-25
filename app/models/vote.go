package models

import "time"

//投票
type Vote struct {
	ID          string `gorm:"type:varchar(20);not null;unique" json:"id"`     //投票主题id
	Initiator   string `gorm:"type:varchar(60);not null" json:"initiator"`     //发起人id
	Title       string `gorm:"type:varchar(120);not null" json:"title"`        //投票标题
	Description string `gorm:"type:varchar(6000);not null" json:"description"` //投票描述，包括规则，富文本
	Img         string `gorm:"type:varchar(300)" json:"img"`                   //图片

	Status      string    `gorm:"type:varchar(10);not null" json:"status"`       //状态，"进行中"、"已结束"
	CreatedTime time.Time `gorm:"type:datetime(6);not null" json:"created_time"` //创建时间
	UpdatedTime time.Time `gorm:"type:datetime(6);not null" json:"updated_time"` //修改时间
	EndTime     time.Time `gorm:"type:datetime(6);not null" json:"end_time"`     //结束时间
}

//投票项
type VoteItem struct {
	VoteID   string `gorm:"type:varchar(20);not null" json:"vote_id"`    //投票主题id
	ItemID   int    `gorm:"type:int;not null" json:"item_id" `           //项id
	ItemName string `gorm:"type:varchar(300);not null" json:"item_name"` //项名字

}

//投票记录
type VoteRecord struct {
	VoteID   string    `gorm:"type:varchar(20);not null" json:"vote_id"`   //投票主题id
	UserID   string    `gorm:"type:varchar(20);not null" json:"user_id"`   //用户id
	Choice   int       `gorm:"type:int;not null" json:"choice"`            //用户选择，对应投票项的项id
	VoteTime time.Time `gorm:"type:datetime(6);not null" json:"vote_time"` //投票时间
}
