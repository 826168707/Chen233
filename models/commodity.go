package models

import "LedgerProject/dao"

type Commodity struct {
	Id int `json:"id"`
	Kind int `json:"kind"`
	Img string `json:"img"`
	Price string `json:"price"`
	Name string `json:"name"`
	Comment string `json:"comment"`
	Path string `json:"path"`
}

func FindKind1() (error, []Commodity) {
	var commodities []Commodity
	err := dao.DB.Where("kind = ?","1").Order("rand()").Limit(6).Find(&commodities).Error
	return err,commodities
}

func FindKind2() (error, []Commodity) {
	var commodities []Commodity
	err := dao.DB.Where("kind = ?","2").Order("rand()").Limit(6).Find(&commodities).Error
	return err,commodities
}

func FindKind3() (error, []Commodity) {
	var commodities []Commodity
	err := dao.DB.Where("kind = ?","3").Order("rand()").Limit(6).Find(&commodities).Error
	return err,commodities
}

func FindKind0() (error, []Commodity) {
	var commodities []Commodity
	err := dao.DB.Where("kind = ?","0").Order("rand()").Limit(6).Find(&commodities).Error
	return err,commodities
}
