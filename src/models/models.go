package models

import (
	"errors"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// 设置数据库路径
	_DB_NAME = "data/goBlog.db"
	// 设置数据库名称
	_SQLITE3_DRIVER = "sqlite3"
)

var (
	engine *xorm.Engine
)

// 分类
type Category struct {
	Id              int64
	Title           string
	Created         time.Time `xorm:"index"`
	Views           int64     `xorm:"index"`
	TopicTime       time.Time `xorm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

// 文章
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Content         string `xorm:"size(5000)"`
	Attachment      string
	Created         time.Time `xorm:"index"`
	Updated         time.Time `xorm:"index"`
	Views           int64     `xorm:"index"`
	Author          string
	ReplyTime       time.Time `xorm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

func RegisterDB() {
	// 检查数据库文件
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	var err error
	engine, err = xorm.NewEngine(_SQLITE3_DRIVER, _DB_NAME)
	if err != nil {
		log.Fatalf("fail to create xorm Engine: %v", err)
	}
	err = engine.Sync(new(Category), new(Topic))
}

func AddCategory(name string) error {
	cate := &Category{
		Title:     name,
		Created:   time.Now(),
		TopicTime: time.Now(),
	}

	// 查询数据
	has, err := engine.Where("title=?", name).Get(cate)
	if has == true {
		return errors.New("category title already exist.")
	}

	// 插入数据
	_, err = engine.Insert(cate)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	cate := &Category{Id: cid}
	_, err = engine.Delete(cate)
	return err
}

func GetAllCategories() ([]*Category, error) {
	cates := make([]*Category, 0)
	err := engine.Find(&cates)
	return cates, err
}
