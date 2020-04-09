package logic

import (
	"LedgerProject/models"
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

func LogicGetHome(email *string) (err error, username,days string, money,avilablemoney int){
	var user models.User
	err,user = models.FindUserByEmail(email)

	username = user.Username
	money = user.Money

	//判断deadline是否为nil----是:days(您还没有设置日期哦!)---否:days为deadline减now的天数
	if user.Deadline == "nil" {
		days = "您还没有设置日期哦!"
		avilablemoney = 0
	}else {
		//计算天数
		CalculateDays(&user.Deadline)
	}


	//计算可用余额

	return
}

func CalculateDays(deadline *string)(days string) {

}

//日期变为时间戳
func DataToTimeStr(deadline *string) (err error) {
	
}