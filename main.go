package main

import (
	"LedgerProject/dao"
	"LedgerProject/models"
	"LedgerProject/routers"
	"fmt"
)

func main() {

	//连接数据库
	err := dao.InitMysql()
	if err != nil {
		fmt.Printf("InitMysql failed , err:%v\n",err)
		return
	}

	defer dao.Close()

	//模型绑定
	dao.DB.AutoMigrate(&models.User{},&models.History{})

	//注册路由
	r := routers.SetupRouter()
	r.Run(":8080")
}
