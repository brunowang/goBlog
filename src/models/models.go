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
	orm *xorm.Engine
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
	orm, err = xorm.NewEngine(_SQLITE3_DRIVER, _DB_NAME)
	if err != nil {
		log.Fatalf("fail to create xorm Engine: %v", err)
	}
	err = orm.Sync(new(Category), new(Topic))
}

func AddCategory(name string) error {
	cate := &Category{
		Title:     name,
		Created:   time.Now(),
		TopicTime: time.Now(),
	}

	// 查询数据
	has, err := orm.Where("title=?", name).Get(cate)
	if has == true {
		return errors.New("category title already exist.")
	}

	// 插入数据
	_, err = orm.Insert(cate)
	return err
}

func DeleteCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	cate := &Category{Id: cid}
	_, err = orm.Delete(cate)
	return err
}

func GetAllCategories() ([]*Category, error) {
	cates := make([]*Category, 0)
	err := orm.Find(&cates)
	return cates, err
}

func AddTopic(title, content string) error {
	topic := &Topic{
		Title:   title,
		Content: content,
		Created: time.Now(),
		Updated: time.Now(),
	}
	_, err := orm.Insert(topic)
	return err
}

func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	topic := new(Topic)
	has, err := orm.Id(tidNum).Get(topic)
	if err != nil {
		return nil, err
	} else if has == false {
		return nil, errors.New("topic id not exist.")
	}

	topic.Views++
	_, err = orm.Id(tidNum).Update(topic)
	return topic, nil
}

func ModifyTopic(tid, title, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	topic := new(Topic)
	has, err := orm.Id(tidNum).Get(topic)
	if has == true {
		topic.Title = title
		topic.Content = content
		topic.Updated = time.Now()
		_, err = orm.Id(tidNum).Update(topic)
	}
	return err
}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	_, err = orm.Delete(topic)
	return err
}

func GetAllTopics(isDesc bool) (topics []*Topic, err error) {
	topics = make([]*Topic, 0)

	if isDesc {
		err = orm.Desc("created").Find(&topics)
	} else {
		err = orm.Find(&topics)
	}
	return topics, err
}
