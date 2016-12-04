package models

import (
	"engine"
	"errors"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

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
	TopicCount      int
	TopicLastUserId int64
}

// 文章
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Category        string
	Labels          string
	Content         string `xorm:"size(5000)"`
	Attachment      string
	Created         time.Time `xorm:"index"`
	Updated         time.Time `xorm:"index"`
	Views           int64     `xorm:"index"`
	Author          string
	ReplyTime       time.Time `xorm:"index"`
	ReplyCount      int
	ReplyLastUserId int64
}

// 评论
type Reply struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func RegisterDB() {
	// 检查数据库文件
	if !engine.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	var err error
	orm, err = xorm.NewEngine(_SQLITE3_DRIVER, _DB_NAME)
	if err != nil {
		log.Fatalf("fail to create xorm Engine: %v", err)
	}
	err = orm.Sync(new(Category), new(Topic), new(Reply))
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

func AddTopic(title, category, label, content, attachment string) error {
	// 处理标签
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	topic := &Topic{
		Title:      title,
		Category:   category,
		Labels:     label,
		Content:    content,
		Attachment: attachment,
		Created:    time.Now(),
		Updated:    time.Now(),
	}
	_, err := orm.Insert(topic)
	if err != nil {
		return err
	}

	// 更新分类统计
	cate := new(Category)
	var topics []*Topic
	topics, err = GetAllTopics(category, label, false)
	if err == nil {
		cate.TopicCount = len(topics)
		_, err = orm.Where("title=?", category).Update(cate)
	}

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

	topic.Labels = strings.Replace(strings.Replace(
		topic.Labels, "#", " ", -1), "$", "", -1)
	return topic, nil
}

func ModifyTopic(tid, title, category, label, content, attachment string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	topic := new(Topic)
	has, err := orm.Id(tidNum).Get(topic)

	var oldCate, oldLabel, oldAttach string
	if has == true {
		oldCate = topic.Category
		oldLabel = topic.Labels
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Category = category
		topic.Labels = label
		topic.Content = content
		topic.Attachment = attachment
		topic.Updated = time.Now()
		_, err = orm.Id(tidNum).Update(topic)
	}
	if err != nil {
		log.Fatal(err)
		return err
	}

	// 更新旧分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		var topics []*Topic
		topics, err = GetAllTopics(oldCate, oldLabel, false)
		if err == nil {
			cate.TopicCount = len(topics)
			_, err = orm.Where("title=?", oldCate).Update(cate)
		}
	}

	// 删除旧的附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	// 更新新分类统计
	cate := new(Category)
	var topics []*Topic
	topics, err = GetAllTopics(category, label, false)
	if err == nil {
		cate.TopicCount = len(topics)
		_, err = orm.Where("title=?", category).Update(cate)
	}

	return err
}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	has, err := orm.Get(topic)

	var oldCate, oldLabel, oldAttach string
	if has == true {
		oldCate = topic.Category
		oldLabel = topic.Labels
		oldAttach = topic.Attachment
		_, err = orm.Delete(topic)
		if err != nil {
			return err
		}
	}

	// 更新分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		var topics []*Topic
		topics, err = GetAllTopics(oldCate, oldLabel, false)
		if err == nil {
			cate.TopicCount = len(topics)
			_, err = orm.Where("title=?", oldCate).Update(cate)
		}
	}

	// 删除附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	return err
}

func GetAllTopics(category, label string, isHomePage bool) (topics []*Topic, err error) {
	topics = make([]*Topic, 0)
	if isHomePage {
		ormSession := orm.Desc("created")
		if len(category) > 0 {
			ormSession = orm.Where("category=?", category)
		}
		if len(label) > 0 {
			ormSession = ormSession.Where("labels like ?", "%$"+label+"#%")
		}
		err = ormSession.Find(&topics)
	} else {
		err = orm.Find(&topics)
	}
	return topics, err
}

func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	reply := &Reply{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	_, err = orm.Insert(reply)
	if err != nil {
		return err
	}

	//更新回复统计
	topic := new(Topic)
	var replies []*Reply
	replies, err = GetAllReplies(tid)
	if err == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount = len(replies)
		_, err = orm.Where("tid=?", tidNum).Update(topic)
	}
	return err
}

func GetAllReplies(tid string) (replies []*Reply, err error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies = make([]*Reply, 0)

	err = orm.Id(tidNum).Find(&replies)
	return replies, err
}

func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}

	reply := &Reply{Id: ridNum}
	has, err := orm.Get(reply)
	var tidNum int64
	if has == true {
		tidNum = reply.Tid
		_, err = orm.Delete(reply)
		if err != nil {
			return err
		}
	}

	replies := make([]*Reply, 0)
	err = orm.Where("tid=?", tidNum).Desc("created").Find(&replies)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	has, err = orm.Get(topic)
	if has == true {
		topic.ReplyTime = replies[0].Created
		topic.ReplyCount = len(replies)
		_, err = orm.Update(topic)
	}
	return err
}
