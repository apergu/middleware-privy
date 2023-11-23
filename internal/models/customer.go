package models

import (
	"middleware/internal/model"
)

type ModelCustomer struct {
	model.Customer
}

//func (filter *ModelCustomer) CreateCustomer() (res model.Customer, err error) {
//	tx := pgxdb.ModelsDB.Begin() // Start a new transaction
//
//	if err = tx.Table("customer").Create(&filter.Customer).Error; err != nil {
//		tx.Rollback()
//		return res, err
//	}
//
//	if err := tx.Commit().Error; err != nil {
//		log.Println(err.Error())
//		return res, err
//	}
//
//	return filter.Customer, nil
//
//}
