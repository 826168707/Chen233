package controller

import (
	"LedgerProject/logic"
	"LedgerProject/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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

//得到主页信息  1.用户名  2.余额 3.可用余额 4.距离截止日期的天数
func GetHome(c *gin.Context)  {
	//根据session中的email从数据库中获取数据
	email := sessions.Default(c).Get("loginuser")
	err := logic.LogicGetHome(&email)
	if err != nil {
		fmt.Printf("GetHome failed err:%v\n",err)
		return
	}
}

//设置金额 截止日期
func SetHome(c *gin.Context)  {
	money := c.PostForm("money")
	deadline := c.PostForm("deadline")	//格式 2006-01-02
	//将日期变为时间戳
	logic.DataToTimeStr(&deadline)
	//根据session中的email从数据库中获取数据
	email := sessions.Default(c).Get("loginuser")
}