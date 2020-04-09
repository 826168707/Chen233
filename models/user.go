package models

import (
	"LedgerProject/dao"
)

//储存用户信息的结构体
type User struct {
	Id 		int	`json:"id"`
	Money	int `json:"money"`
	Deadline string `json:"deadline" gorm:"default:'nil'"`
	Username string	`json:"username"`
	Email 	 string `json:"email"`
	Password string `json:"password"`
	Placeholder string		//预留数据项
}

func CreateNewUser(username, email, password string) (err error) {
	newUser := User{
		Username:username,
		Email:email,
		Password:password,
	}
	err = dao.DB.Create(&newUser).Error
	return
}

func FindUserByEmail(email *string) (err error,user User) {

	err = dao.DB.Where("email=?",*email).First(&user).Error

	return
}

func FindUserByEmailAndPassword(email, password *string) (err error, user User) {

	err = dao.DB.Where("email=? AND password=?",*email,*password).First(&user).Error
	return
}