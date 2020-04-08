package controller

import (
	"LedgerProject/logic"
	"LedgerProject/models"
	"fmt"
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
		//变为用户登录态
		//------

	}
}