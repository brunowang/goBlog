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

func AddTopic(uid int64, title, category, label, content, attachment string) error {
	// 处理标签
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	topic := &db.Topic{
		Uid:        uid,
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
	topics, err = GetAllTopics(uid, category, "", false)
	if err == nil {
		cate.TopicCount = len(topics)
		_, err = db.GetOrm().Where("title=?", category).Cols("topic_count").Update(cate)
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

func ModifyTopic(uid int64, tid, title, category, label, content, attachment string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	topic := new(db.Topic)
	has, err := db.GetOrm().Id(tidNum).Get(topic)

	var oldCate, oldAttach string
	if has == true {
		oldCate = topic.Category
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
		topics, err = GetAllTopics(uid, oldCate, "", false)
		if err == nil {
			cate.TopicCount = len(topics)
			_, err = db.GetOrm().Where("title=?", oldCate).Cols("topic_count").Update(cate)
		}
	}

	// 删除旧的附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	// 更新新分类统计
	cate := new(db.Category)
	var topics []*db.Topic
	topics, err = GetAllTopics(uid, category, "", false)
	if err == nil {
		cate.TopicCount = len(topics)
		_, err = db.GetOrm().Where("title=?", category).Cols("topic_count").Update(cate)
	}

	return err
}

func DeleteTopic(uid int64, tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	topic := &db.Topic{Id: tidNum}
	has, err := db.GetOrm().Get(topic)

	var oldCate, oldAttach string
	if has == true {
		oldCate = topic.Category
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
		topics, err = GetAllTopics(uid, oldCate, "", false)
		if err == nil {
			cate.TopicCount = len(topics)
			_, err = db.GetOrm().Where("title=?", oldCate).Cols("topic_count").Update(cate)
		}
	}

	// 删除附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	return err
}

func GetAllTopics(uid int64, category, label string, isHomePage bool) (topics []*db.Topic, err error) {
	topics = make([]*db.Topic, 0)
	if uid == -1 {
		return topics, nil
	}
	ormSession := db.GetOrm().Where("uid=?", uid)
	if isHomePage {
		ormSession = ormSession.Desc("created")
	}
	if len(category) > 0 {
		ormSession = ormSession.Where("category=?", category)
	}
	if len(label) > 0 {
		ormSession = ormSession.Where("labels like ?", "%$"+label+"#%")
	}
	err = ormSession.Find(&topics)
	return topics, err
}
