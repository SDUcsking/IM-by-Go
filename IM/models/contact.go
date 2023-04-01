package models

type Contact struct {
	Id       int `gorm:"primary_key;AUTO_INCREMENT"`
	OwnerId  int
	TargetId int
	Type     int
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}
