package logic

import (
	"LedgerProject/models"
	"fmt"
)

//不存在返回true  存在返回false
func EmailCheck(email *string) bool {
	//检查数据库中是否存在相同email
	err,user := models.FindUserByEmail(email)
	if err != nil {
		panic(err)
	}
	if user.Email != "" {
		return false
	}
	return true
}

//存在用户返回true  不存在返回false
func UserCheck(email, password *string) bool {
	err,user := models.FindUserByEmailAndPassword(email,password)
	if err != nil {
		fmt.Printf("FindUserByEmailAndPassword failed,err:%v",err)
		return false
	}
	if user.Email != "" {
		return true
	}
	return false
}