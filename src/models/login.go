package models

import (
	"db"
)

func CheckAccount(name, pwd string) bool {
	acc := new(db.Account)
	has, _ := db.GetOrm().Where("name=? and pwd=?", name, pwd).Get(acc)
	return has
}

func GetAccount(name, pwd string) *db.Account {
	acc := new(db.Account)
	has, _ := db.GetOrm().Where("name=? and pwd=?", name, pwd).Get(acc)
	if has {
		return acc
	}
	return nil
}
