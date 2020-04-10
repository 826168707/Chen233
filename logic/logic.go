package logic

import (
	"LedgerProject/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

//不存在返回true  存在返回false
func EmailCheck(email *string) bool {
	//检查数据库中是否存在相同email
	err,_:= models.FindUserByEmail(email)
	if err != nil {		//没有找到相同数据
		return true
	}else {
		return false
	}
}

//存在用户返回true  不存在返回false
func UserCheck(email, password *string) bool {
	err,_ := models.FindUserByEmailAndPassword(email,password)
	if err != nil {		//没找到相同数据
		return false
	}else {
		return true
	}

}


func LogicGetHome(email string) (err error, username,days string, money,usefulMoney int){
	var user models.User
	err,user = models.FindUserByEmail(&email)

	username = user.Username
	money = user.Money
	

	//判断deadline是否为nil----是:days(您还没有设置日期哦!)---否:days为deadline减now的天数
	if user.Deadline == "nil" {
		days = "您还没有设置日期哦!"
		usefulMoney = 0
	}else {
		//检测日期是否过期
		if !IsExpired(user.Deadline){ //过期的情况
			days = "设置的日期已经过期啦,请重新设置本月计划吧!"
			usefulMoney = money
		}else {		//未过期
			//计算天数
			days_Int := CalculateDays(user.Deadline)
			days = strconv.Itoa(days_Int)
			//计算可用余额
			usefulMoney = CalculateUsefulMoney(days_Int,money,user.Dailyexpenses)

		}
	}
	return
}

//检测是否过期  true表未过期
func IsExpired(deadlline string) bool {
	t2 := StringToTime(deadlline)
	t1 := time.Now()
	return t1.Before(t2)
}


//计算可用余额  可用余额=余额-天数*每日固定支出
func CalculateUsefulMoney(days,money,dailyexpenses int)(usefulMoney int)  {
	usefulMoney = money - days * dailyexpenses
	return 
}


func CalculateDays(deadline string)(days int) {
	t1 := time.Now()
	t2 := StringToTime(deadline)

	days = int(t2.Sub(t1).Hours()/24)

	return
}

//从session获取email并进行类型修正
func GetEmailFromSession(c *gin.Context)(email string){
	interEmail := sessions.Default(c).Get("loginuser")
	email = interEmail.(string)
	return
}
//日期变为时间戳
func DataToTimeStr(deadline *string) (timeSta int64,err error) {
	timeLayout := "2006-01-02"
	loc,err := time.LoadLocation("Asia/Shanghai")
	theTime,err := time.ParseInLocation(timeLayout,*deadline,loc)
	timeSta = theTime.Unix()
	return
}

//将string变为time类型
func StringToTime(str string) (theTime time.Time) {
	timeLayout := "2006-01-02"
	loc,_ := time.LoadLocation("Asia/Shanghai")
	theTime,_ = time.ParseInLocation(timeLayout,str,loc)
	return
}


//更新对应用户的金额,日期,日常固定支出
func UpdateCount(email, deadline string, money ,dailyexpenses int) (err error) {

	//更新
	err = models.UpdateMoneyAndDeadline(email,money,dailyexpenses,deadline)
	return
}