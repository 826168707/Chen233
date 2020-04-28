package routers

import (
	"LedgerProject/controller"
	"LedgerProject/logic"
	"github.com/gin-gonic/gin"
)


func SetupRouter() *gin.Engine {

	r := gin.Default()


	r.Use(controller.Cors())
	//r.Static()	导入静态文件
	//r.LoadHTMLGlob()	模板

	//store := cookie.NewStore([]byte("loginuser"))
	//r.Use(sessions.Sessions("session",store))

	//注册登录相关路由组
	v1Group := r.Group("sign")
	{
		//注册
		v1Group.PUT("/up",controller.UserRegistered)
		//发送验证码
		v1Group.POST("/up",controller.SendEmail)
		//登录
		v1Group.POST("",controller.UserLogin)
	}

	//主页路由组
	v2Group := r.Group("home")
	{
		//登录主页后页面获取信息
		v2Group.GET("",logic.JWTAuthMiddleware(),controller.GetHome)

		//设置金额 截止日期  日常固定支出
		v2Group.PUT("",logic.JWTAuthMiddleware(),controller.SetHome)

		//退出登录
		v2Group.POST("/out",logic.JWTAuthMiddleware(),controller.UserSignOut)

	}

	//支出 收录 路由组
	v3Group := r.Group("set")
	{
		//想要添加特殊支出
		v3Group.POST("/cost",logic.JWTAuthMiddleware(),controller.WantCost)
		//确认支出
		v3Group.PUT("/cost",logic.JWTAuthMiddleware(),controller.AddCost)
		//添加收入
		v3Group.PUT("/income",logic.JWTAuthMiddleware(),controller.AddIncome)
	}

	//历史记录路由组
	v4Group := r.Group("history")
	{
		//支出历史记录
		v4Group.GET("/cost",logic.JWTAuthMiddleware(),controller.CostHistory)
		//修改历史记录
		v4Group.PUT("",logic.JWTAuthMiddleware(),controller.UpdateHistory)
		//收入历史记录
		v4Group.GET("/income",logic.JWTAuthMiddleware(),controller.IncomeHistory)
		//删除历史记录
		v4Group.DELETE("",logic.JWTAuthMiddleware(),controller.DeleteHistory)
	}

	//推荐路由
	r.GET("/recommend",logic.JWTAuthMiddleware(),controller.Recommend)


	return r
}



