package routers

import (
	"LedgerProject/controller"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)


func SetupRouter() *gin.Engine {
	r := gin.Default()

	//r.Static()	导入静态文件
	//r.LoadHTMLGlob()	模板

	store := cookie.NewStore([]byte("loginuser"))
	r.Use(sessions.Sessions("session",store))

	//注册登录相关路由组
	v1Group := r.Group("sign")
	{
		//注册
		v1Group.POST("/up",controller.UserRegistered)
		//登录
		v1Group.POST("/",controller.UserLogin)
	}

	//主页路由组
	v2Group := r.Group("home")
	{
		//登录主页后页面获取信息
		v2Group.GET("/",controller.GetHome)

		//设置金额 截止日期
		v2Group.PUT("/",controller.SetHome)

		//退出登录
		v2Group.POST("/out",controller.UserSignOut)
	}






	return r
}

