package models

import (
	"db"
	"errors"
	"strconv"
	"time"
)

func AddCategory(name string) error {
	cate := new(db.Category)

	// 查询数据
	has, err := db.GetOrm().Where("title=?", name).Get(cate)
	if has == true {
		return errors.New("category title already exist.")
	}

	// 插入数据
	_, err = db.GetOrm().Insert(&db.Category{
		Title:     name,
		Created:   time.Now(),
		TopicTime: time.Now(),
	})
	return err
}

func DeleteCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	cate := &db.Category{Id: cid}
	_, err = db.GetOrm().Delete(cate)
	return err
}

func GetAllCategories() ([]*db.Category, error) {
	cates := make([]*db.Category, 0)
	err := db.GetOrm().Find(&cates)
	return cates, err
}
