package models

type GroupBasic struct {
	Id      int `gorm:"primary_key;AUTO_INCREMENT"`
	Name    string
	OwnerId int
	Icon    string
	Desc    string
	Type    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
