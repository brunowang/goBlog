package models

import (
	"db"
	"errors"
)

func AddAccount(name string, pwd string) error {
	acc := &db.Account{
		Name: name,
		Pwd:  pwd,
	}

	has, _ := db.GetOrm().Where("name=?", name).Get(acc)
	if has == true {
		return errors.New("account name duplicate.")
	}

	_, err := db.GetOrm().Insert(acc)
	return err
}
