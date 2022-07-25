package models

import "time"

type User struct {
	ID           string    `gorm:"type:varchar(20);not null;unique" `             //用户ID
	Name         string    `gorm:"type:varchar(30)" json:"name"`                  //姓名
	Phone        string    `gorm:"type:varchar(11);not null;unique" json:"phone"` //电话
	HeadImg      string    `gorm:"type:varchar(300)" json:"head_img"`             //头像
	IDCard       string    `gorm:"type:varchar(18)" json:"id_card"`               //身份证号
	OccupiedArea float64   `gorm:"type:decimal(6,1)" json:"occupied_area"`        //占有面积
	Proof        string    `gorm:"type:varchar(300)" json:"proof"`                //房产证照片
	PassWord     string    `gorm:"type:varchar(50);not null" json:"pass_word"`    //密码
	Auth         int       `gorm:"type:int;not null" json:"auth"`                 //用户权限，1未提交审核2已提交审核3审核通过4驳回
	SignupTime   time.Time `gorm:"type:datetime(6)" json:"signup_time"`           //注册时间
	SubmitTime   time.Time `gorm:"type:datetime(6)" json:"submit_time"`           //审核信息提交时间
	AdoptTime    time.Time `gorm:"type:datetime(6)" json:"adopt_time"`            //审核通过时间
}
