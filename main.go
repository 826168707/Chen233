package main

import (
	"LedgerProject/dao"
	"LedgerProject/log"
	"LedgerProject/models"
	"LedgerProject/routers"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	//初始化日志库
	log.SetLogs(zap.DebugLevel,log.LOGFORMAT_CONSOLE,"./log/server.log")
	//连接数据库
	if err := dao.InitMysql();err != nil {
		fmt.Printf("InitMysql failed , err:%v\n",err)
		return
	}
	if err := dao.InitRedis();err != nil {
		fmt.Printf("InitRedis failed ,err:%v\n",err)
	}

	defer dao.Close()
	defer dao.Rclose()

	//模型绑定
	dao.DB.AutoMigrate(&models.User{},&models.History{},&models.Commodity{})

	//注册路由
	r := routers.SetupRouter()
	r.Run(":7777")
}
