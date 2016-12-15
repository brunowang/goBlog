package models

import (
	"db"
	"errors"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func AddTopic(title, category, label, content, attachment string) error {
	// 处理标签
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	topic := &db.Topic{
		Title:      title,
		Category:   category,
		Labels:     label,
		Content:    content,
		Attachment: attachment,
		Created:    time.Now(),
		Updated:    time.Now(),
	}
	_, err := db.GetOrm().Insert(topic)
	if err != nil {
		return err
	}

	// 更新分类统计
	cate := new(db.Category)
	var topics []*db.Topic
	topics, err = GetAllTopics(category, label, false)
	if err == nil {
		cate.TopicCount = len(topics)
		_, err = db.GetOrm().Where("title=?", category).Update(cate)
	}

	return err
}

func GetTopic(tid string) (*db.Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	topic := new(db.Topic)
	has, err := db.GetOrm().Id(tidNum).Get(topic)
	if err != nil {
		return nil, err
	} else if has == false {
		return nil, errors.New("topic id not exist.")
	}

	topic.Views++
	_, err = db.GetOrm().Id(tidNum).Update(topic)

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

	topic := new(db.Topic)
	has, err := db.GetOrm().Id(tidNum).Get(topic)

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
		_, err = db.GetOrm().Id(tidNum).Update(topic)
	}
	if err != nil {
		log.Fatal(err)
		return err
	}

	// 更新旧分类统计
	if len(oldCate) > 0 {
		cate := new(db.Category)
		var topics []*db.Topic
		topics, err = GetAllTopics(oldCate, oldLabel, false)
		if err == nil {
			cate.TopicCount = len(topics)
			_, err = db.GetOrm().Where("title=?", oldCate).Update(cate)
		}
	}

	// 删除旧的附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	// 更新新分类统计
	cate := new(db.Category)
	var topics []*db.Topic
	topics, err = GetAllTopics(category, label, false)
	if err == nil {
		cate.TopicCount = len(topics)
		_, err = db.GetOrm().Where("title=?", category).Update(cate)
	}

	return err
}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	topic := &db.Topic{Id: tidNum}
	has, err := db.GetOrm().Get(topic)

	var oldCate, oldLabel, oldAttach string
	if has == true {
		oldCate = topic.Category
		oldLabel = topic.Labels
		oldAttach = topic.Attachment
		_, err = db.GetOrm().Delete(topic)
		if err != nil {
			return err
		}
	}

	// 更新分类统计
	if len(oldCate) > 0 {
		cate := new(db.Category)
		var topics []*db.Topic
		topics, err = GetAllTopics(oldCate, oldLabel, false)
		if err == nil {
			cate.TopicCount = len(topics)
			_, err = db.GetOrm().Where("title=?", oldCate).Update(cate)
		}
	}

	// 删除附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	return err
}

func GetAllTopics(category, label string, isHomePage bool) (topics []*db.Topic, err error) {
	topics = make([]*db.Topic, 0)
	if isHomePage {
		ormSession := db.GetOrm().Desc("created")
		if len(category) > 0 {
			ormSession = db.GetOrm().Where("category=?", category)
		}
		if len(label) > 0 {
			ormSession = ormSession.Where("labels like ?", "%$"+label+"#%")
		}
		err = ormSession.Find(&topics)
	} else {
		err = db.GetOrm().Find(&topics)
	}
	return topics, err
}
