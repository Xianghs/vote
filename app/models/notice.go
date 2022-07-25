package models

import "time"

//公告
type Notice struct {
	ID          string    `gorm:"type:varchar(20);not null" json:"id"`            //公告id
	Initiator   string    `gorm:"type:varchar(60);not null" json:"initiator"`     //发起人名称
	Title       string    `gorm:"type:varchar(120);not null" json:"title"`        //公告标题
	IsTop       int       `gorm:"type:int;not null" json:"is_top"`                //是否置顶，1不置顶2置顶
	Img         string    `gorm:"type:varchar(300)" json:"img"`                   //图
	Abstract    string    `gorm:"type:varchar(300)" json:"abstract"`              //摘要
	Description string    `gorm:"type:varchar(6000);not null" json:"description"` //公告内容
	CreatedTime time.Time `gorm:"type:datetime(6);not null" json:"created_time"`  //创建时间
	UpdatedTime time.Time `gorm:"type:datetime(6);not null" json:"updated_time"`  //修改时间
}
