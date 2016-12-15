package models

import (
	"db"
)

func CheckAccount(name string, pwd string) bool {
	acc := new(db.Account)
	has, _ := db.GetOrm().Where("name=? and pwd=?", name, pwd).Get(acc)
	return has
}
