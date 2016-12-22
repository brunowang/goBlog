package models

import (
	"db"
	"strconv"
	"time"
)

func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	reply := &db.Reply{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	_, err = db.GetOrm().Insert(reply)
	if err != nil {
		return err
	}

	//更新回复统计
	topic := new(db.Topic)
	var replies []*db.Reply
	replies, err = GetAllReplies(tid)
	if err == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount = len(replies)
		_, err = db.GetOrm().Id(tidNum).Update(topic)
	}
	return err
}

func GetAllReplies(tid string) (replies []*db.Reply, err error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies = make([]*db.Reply, 0)

	err = db.GetOrm().Where("tid=?", tidNum).Find(&replies)
	return replies, err
}

func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}

	reply := &db.Reply{Id: ridNum}
	has, err := db.GetOrm().Get(reply)
	var tidNum int64
	if has == true {
		tidNum = reply.Tid
		_, err = db.GetOrm().Delete(reply)
		if err != nil {
			return err
		}
	}

	replies := make([]*db.Reply, 0)
	err = db.GetOrm().Where("tid=?", tidNum).Desc("created").Find(&replies)
	if err != nil {
		return err
	}

	topic := &db.Topic{}
	has, err = db.GetOrm().Get(topic)
	if has == true {
		if len(replies) != 0 {
			topic.ReplyTime = replies[0].Created
		}
		topic.ReplyCount = len(replies)
		_, err = db.GetOrm().Id(tidNum).Update(topic)
	}
	return err
}
