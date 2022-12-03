package model

type FriendLink struct {
	Universal
	Name    string `gorm:"type:varchar(50)" json:"name"`
	Avatar  string `gorm:"type:varchar(255)" json:"avatar"`
	Address string `gorm:"type:varchar(255)" json:"address"`
	Intro   string `gorm:"type:varchar(255)" json:"intro"`
}
