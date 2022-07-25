package models

//管理员账号
type Admin struct {
	ID        string `gorm:"type:varchar(20);not null" json:"id"`         //管理员id
	AdminName string `gorm:"type:varchar(60);not null" json:"admin_name"` //管理员名称
	Level     int    `gorm:"type:int;not null" json:"level"`              //管理员级别1超级管理员2普通管理员
	Password  string `gorm:"type:varchar(50);not null" json:"password"`   //管理员密码
}
