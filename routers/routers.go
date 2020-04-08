package routers

import (
	"LedgerProject/controller"
	"github.com/gin-gonic/gin"

)


func SetupRouter() *gin.Engine {
	r := gin.Default()

	//r.Static()	导入静态文件
	//r.LoadHTMLGlob()	模板

	//注册登录相关路由组
	v1Group := r.Group("sign")
	{
		//注册
		v1Group.POST("/up",controller.UserRegistered)
		//登录
		v1Group.POST("/in",controller.UserLogin)
	}





	return r
}

