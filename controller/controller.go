package controller

import (
	"LedgerProject/logic"
	"LedgerProject/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

//用户注册
func UserRegistered(c *gin.Context)  {
	//前端页面填写 用户名 邮箱 密码 点击提交 请求发送到这里

	//1.把请求中数据提取
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	captcha := c.PostForm("captcha")
	//检测验证码是否正确
	if !logic.CaptchaCheck(captcha){	//错误
		c.JSON(http.StatusOK,gin.H{
			"failed":"验证码错误!",
		})
		return
	}

	//不存在相同邮箱，注册用户
	err := models.CreateNewUser(username,email,password)
	if err != nil {
		fmt.Printf("CreateNewUser failed, err:%v",err)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"注册完成！",
	})


}

//发送验证码
func SendEmail(c *gin.Context)  {
	email := c.PostForm("email")

	//.判断该邮箱是否已经被使用
	if logic.EmailCheck(&email) == false {	//以被使用
		c.JSON(http.StatusOK,gin.H{
			"failed":"该邮箱已经被使用！",
		})
		return
	}

	if !logic.SendEmail(email){	//false
		c.JSON(http.StatusOK,gin.H{
			"failed":"还没有填写邮箱哦!",
		})
		return
	}


	c.JSON(http.StatusOK,gin.H{
		"message":"验证码已发送",
	})

}

//用户登录
func UserLogin(c *gin.Context)  {
	//前段填写 邮箱 密码 发送到此
	email := c.PostForm("email")
	password := c.PostForm("password")

	//检查数据库中是否存在该用户
	if logic.UserCheck(&email,&password) == false{		//不存在该用户
		c.JSON(http.StatusOK,gin.H{
			"message":"邮箱或密码错误!",
		})
	}else { 		//成功登录
		//添加session
		//session := sessions.Default(c)
		//session.Set("loginuser",email)
		//session.Save()
		//fmt.Printf(session.Get("loginuser").(string))

		//生成Token
		token,_ := logic.GenToken(email)

		c.JSON(http.StatusOK,gin.H{
			"message":"登录成功!",
			"token":token,
		})
	}
}

//获取session
func GetSession(c *gin.Context) bool {
	session := sessions.Default(c)
	loginuser := session.Get("loginuser")
	if loginuser != nil {
		return true
	}else {
		return false
	}
}

//退出登录
func UserSignOut(c *gin.Context)  {
	////删除session中的数据
	//session := sessions.Default(c)
	//session.Delete("loginuser")
	//session.Save()
	//转到注册登录界面
	c.Redirect(http.StatusMovedPermanently,"/sign")
}

//得到主页信息  1.用户名  2.余额 3.可用余额 4.距离截止日期的天数 + 数据可视化
func GetHome(c *gin.Context)  {
	////根据session中的email从数据库中获取数据
	//email := logic.GetEmailFromSession(c)
	interEmail, _ := c.Get("email")
	email := interEmail.(string)

	err,username,days,money,usefulMoney := logic.LogicGetHome(email)
	if err != nil {
		fmt.Printf("GetHome failed err:%v\n",err)
		return
	}
	//可视数据
	err,visualData := logic.VisualData(email)
	if err != nil {
		fmt.Printf("VisualData failed, err:%v\n",err)
		return
	}


	c.JSON(http.StatusOK,gin.H{
		"username":username,
		"money":money,
		"useful_money":usefulMoney,
		"days":days,	//days为string类型
		"visual_data":visualData,
	})
}

//设置金额 截止日期  每日固定支出
func SetHome(c *gin.Context)  {

	var(
		moneyStr = c.PostForm("money")
		deadline = c.PostForm("deadline")	//格式 2006-01-02
		dailyExpensesStr = c.PostForm("daily_expenses")
		interEmail, _ = c.Get("email")
		email = interEmail.(string)
		money,_ = strconv.Atoi(moneyStr)
		dailyExpenses,_ = strconv.Atoi(dailyExpensesStr)
	)

	err := logic.UpdateCount(email,deadline,money,dailyExpenses)
	if err != nil {
		fmt.Printf("UpdateCount failed ,err:%v\n",err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"设置完成!",
	})

}


//想要添加特殊支出
func WantCost(c *gin.Context)  {
	//获取花费金额
	cost, _ := strconv.Atoi(c.PostForm("cost"))
	interEmail, _ := c.Get("email")
	email := interEmail.(string)

	//提示模块
	err,remainMoney := logic.CostTip(email,cost)
	if err != nil {
		fmt.Printf("CostTip failed, err:%v\n,",err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"remain_money":remainMoney,
	})

}

//确认添加特殊支出
func AddCost(c *gin.Context)  {

	//获取种类,注释,花费金额
	kind, _ := strconv.Atoi(c.PostForm("kind"))
	comment := c.PostForm("comment")
	cost, _ := strconv.Atoi(c.PostForm("cost"))
	interEmail, _ := c.Get("email")
	email := interEmail.(string)

	//修改用户金额
	err := models.UpdateMoneyByEmail(email,cost)
	if err != nil {
		fmt.Printf("UpdateMoneyByEmail failed , err:%v\n",err)
		return
	}

	err = logic.AddHistory(email,kind,cost,comment)
	if err != nil {
		fmt.Printf("AddHistory failed ,err :%v\n",err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"完成喽!",
	})

}

//添加收入
func AddIncome(c *gin.Context) {
	//获取 注释 收入金额
	comment := c.PostForm("comment")
	income, _ := strconv.Atoi(c.PostForm("income"))
	income = 0 - income		//变换符号
	interEmail, _ := c.Get("email")
	email := interEmail.(string)
	//修改用户金额
	err := models.UpdateMoneyByEmail(email,income)
	if err != nil {
		fmt.Printf("UpdateMoneyByEmail failed , err:%v\n",err)
		return
	}
	//添加历史记录
	err = logic.AddHistory(email,0,income,comment)
	if err != nil {
		fmt.Printf("AddHistory failed ,err :%v\n",err)
		return
	}

	//提示模块
	err,remainMoney := logic.CostTip(email,income)
	if err != nil {
		fmt.Printf("CostTip failed, err:%v\n,",err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"remain_money":remainMoney,
	})



}

//支出历史记录
func CostHistory(c *gin.Context)  {
	interEmail, _ := c.Get("email")
	email := interEmail.(string)
	//根据email获取所有支出记录
    err,histories := logic.GetCostHistory(email)
    if err != nil {
		fmt.Printf("GetCostHistory failed err:%v\n",err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"histories":histories,
	})

}

//收入历史记录
func IncomeHistory(c *gin.Context)  {
	interEmail, _ := c.Get("email")
	email := interEmail.(string)

	err,histories := logic.GetIncomeHistory(email)
	if err != nil {
		fmt.Printf("GetIncomeHistory failed, err:%v\n",err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"histories":histories,
	})
}

//推荐模块
func Recommend(c *gin.Context)  {
	interEmail, _ := c.Get("email")
	email := interEmail.(string)

	err,commodities := logic.GetRecommend(email)
	if err != nil {
		fmt.Printf("GetRecommend failed , err:%v\n",err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"commodities":commodities,
	})
}

//编辑历史记录
func UpdateHistory(c *gin.Context)  {
	var(
		id,_= strconv.Atoi(c.PostForm("id"))
		kind,_ = strconv.Atoi(c.PostForm("kind"))
		money,_ = strconv.Atoi(c.PostForm("money"))
		date = c.PostForm("date")
		comment = c.PostForm("comment")
	)

	err := models.UpdateHistoriesById(date,comment,id,kind,money)
	if err != nil {
		fmt.Printf("UpdateCostHistoriesById failed,err:%v\n",err)
		c.JSON(http.StatusOK,gin.H{
			"failed":"编辑失败!",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"编辑完成!",
	})
}

//删除历史记录
func DeleteHistory(c *gin.Context)  {
	id,_:= strconv.Atoi(c.PostForm("id"))
	income, _ := strconv.Atoi(c.PostForm("money"))
	income = 0 - income		//变换符号
	interEmail, _ := c.Get("email")
	email := interEmail.(string)

	err := models.DeleteHistoriesById(id)
	if err != nil {
		fmt.Printf("DeleteHistoriesById failed,err:%v\n",err)
		c.JSON(http.StatusOK,gin.H{
			"failed":"删除失败!",
		})
		return
	}

	err = models.UpdateMoneyByEmail(email,income)
	if err != nil {
		fmt.Printf("UpdateMoneyByEmail failed , err:%v\n",err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"删除成功!",
	})

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method		//请求方法
		origin := c.Request.Header.Get("Origin")		//请求头部
		var headerKeys []string								// 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")		// 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")		//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//				允许跨域设置																										可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")		// 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")		// 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")		//	跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")		// 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next()		//	处理请求
	}
}


