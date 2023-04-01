package main

import (
	"IM/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/ginChat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Contact{})
	db.AutoMigrate(&models.GroupBasic{})

	//user := &models.UserBasic{}
	//user.Name = "jack"
	//
	//user.PassWord = "123"
	//db.Create(user)
	//fmt.Println(db.First(user, 1))
	//db.Model(user).Update("PassWord", "1234")
}
