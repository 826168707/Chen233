package main

import (
	"LedgerProject/dao"
	"LedgerProject/models"
	"LedgerProject/routers"
)

func main() {

	//连接数据库
	err := dao.InitMysql()
	if err != nil {
		panic(err)
	}
	defer dao.Close()

	//模型绑定
	dao.DB.AutoMigrate(&models.User{})

	//注册路由
	r := routers.SetupRouter()
	r.Run(":8080")
}
