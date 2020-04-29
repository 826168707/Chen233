package models

import (
	"LedgerProject/dao"
)

type History struct {
	Id int	`json:"id"`
	Kind int `json:"kind"`  //1.餐饮 2.购物 3.交通 4.教育 5.娱乐 6.服饰 7.恋爱 8.水果 9.运动 11.日常支出  0.收入
	Money int `json:"money"`
	Date string `json:"date"`
	Comment string `json:"comment"`
	Email string `json:"email"`
	Placeholder string		//预留数据项
}

func AddOneHistory(email string,kind int,money int,comment string,date string)(err error)  {
	newHistory := History{
		Kind:        kind,
		Money:       money,
		Date:        date,
		Comment:     comment,
		Email:       email,
		Placeholder: "",
	}

	err = dao.DB.Create(&newHistory).Error
	return
}

func FindCostHistoriesByEmail(email string) (err error,histories []History) {

	err = dao.DB.Where("email = ? AND kind <> ?",email,"0").Find(&histories).Error
	return
}


func FindIncomeHistoriesByEmail(email string) (err error,histories []History) {

	err = dao.DB.Where("email = ? AND kind = ?",email,"0").Find(&histories).Error
	return
}

func FindIdFromHistory(email string,kind int,money int,comment string)(err error,history History)  {
	err = dao.DB.Where("eamil = ? AND kind = ? AND money = ? AND date = ?",email,kind,money,comment).First(&history).Error
	return
}

func UpdateHistoriesById(date,comment string ,id,kind,money int) (err error) {
	var history History
	history.Id = id
	err = dao.DB.Model(&history).Update(map[string]interface{}{"kind":kind,"money":money,"date":date,"comment":comment}).Error
	return
}

func DeleteHistoriesById(id int) (err error) {
	var history History
	err = dao.DB.Where("id = ?",id).Delete(&history).Error
	return
}

