package models

import (
	"LedgerProject/dao"

)

type History struct {
	Kind int `json:"kind"`  //1.餐饮 2.购物 3.交通 4.教育 5.娱乐  0.收入
	Money int `json:"money"`
	Data string `json:"data"`
	Comment string `json:"comment"`
	Email string `json:"email"`
	Placeholder string		//预留数据项
}

func AddOneHistory(email string,kind int,money int,comment string,data string)(err error)  {
	newHitory := History{
		Kind:        kind,
		Money:       money,
		Data:        data,
		Comment:     comment,
		Email:       email,
		Placeholder: "",
	}

	err = dao.DB.Create(&newHitory).Error
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