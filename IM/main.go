package main

import (
	"IM/router"
	"IM/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySql()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8080")

}
