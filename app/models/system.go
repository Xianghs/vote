package models

//系统设置
type SystemSet struct {
	ID             int     `gorm:"type:int" json:"id"`
	TotalHousehold int     `gorm:"type:int" json:"total_household"`      //总户数
	TotalArea      float64 `gorm:"type:decimal(10,1)" json:"total_area"` //小区房屋总面积
}
