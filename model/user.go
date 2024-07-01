package model

type Interest struct {
	Name      string `json:"name"`
	OtherInfo string `json:"other_info"`
}

type User struct {
	BaseModel
	Name      string      `gorm:"column:name" json:"name"`
	Phone     string      `gorm:"column:phone" json:"phone"`
	Password  string      `gorm:"password" json:"password"`
	Interests []*Interest `gorm:"serializer:json" json:"interests"`
	CommonTimestampsField
}
