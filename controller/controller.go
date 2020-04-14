package controller

import (
	"LedgerProject/logic"
	"LedgerProject/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//用户注册
func UserRegistered(c *gin.Context)  {
	//前端页面填写 用户名 邮箱 密码 点击提交 请求发送到这里

	//1.把请求中数据提取
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	//2.判断该邮箱是否已经被使用
	if logic.EmailCheck(&email) == false {	//以被使用
		c.JSON(http.StatusOK,gin.H{
			"message":"该邮箱已经被使用！",
		})
	}else {		//不存在相同邮箱，注册用户
		err := models.CreateNewUser(username,email,password)
		if err != nil {
			fmt.Printf("CreateNewUser failed, err:%v",err)
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"message":"注册完成！",
		})
		//返回注册登录页
		c.Redirect(http.StatusMovedPermanently,"/sign")
	}
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
		session := sessions.Default(c)
		session.Set("loginuser",email)
		session.Save()

		//转到主页
		c.Redirect(http.StatusMovedPermanently,"/home")
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
	//删除session中的数据
	session := sessions.Default(c)
	session.Delete("loginuser")
	session.Save()
	//转到注册登录界面
	c.Redirect(http.StatusMovedPermanently,"/sign")
}

//得到主页信息  1.用户名  2.余额 3.可用余额 4.距离截止日期的天数 + 数据可视化
func GetHome(c *gin.Context)  {
	//根据session中的email从数据库中获取数据
	email := logic.GetEmailFromSession(c)

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
		email = logic.GetEmailFromSession(c)
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

	c.Redirect(http.StatusMovedPermanently,"/home")
}


//想要添加特殊支出
func WantCost(c *gin.Context)  {
	//获取花费金额
	cost, _ := strconv.Atoi(c.PostForm("cost"))
	email := logic.GetEmailFromSession(c)

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
	email := logic.GetEmailFromSession(c)

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

	c.Redirect(http.StatusMovedPermanently,"/home")
}

//添加收入
func AddIncome(c *gin.Context) {
	//获取 注释 收入金额
	comment := c.PostForm("comment")
	income, _ := strconv.Atoi(c.PostForm("income"))
	income = 0 - income		//变换符号
	email := logic.GetEmailFromSession(c)

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

	c.Redirect(http.StatusMovedPermanently,"/home")

}

//支出历史记录
func CostHistory(c *gin.Context)  {
	email := logic.GetEmailFromSession(c)

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
	email := logic.GetEmailFromSession(c)

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
	email := logic.GetEmailFromSession(c)

	err,commodities := logic.GetRecommend(email)
	if err != nil {
		fmt.Printf("GetRecommend failed , err:%v\n",err)
	}

	c.JSON(http.StatusOK,gin.H{
		"commodities":commodities,
	})
}