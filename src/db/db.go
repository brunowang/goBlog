package db

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	// 设置Mysql数据源
	_MYSQL_DB_NAME = "brunowang:111111@(123.56.29.218:3306)/goBlog"
	// 设置Mysql驱动名称
	_MYSQL_DRIVER = "mysql"
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

// 账号
type Account struct {
	Id   int64
	Name string
	Pwd  string
}

func RegisterDB() {
	var err error
	orm, err = xorm.NewEngine(_MYSQL_DRIVER, _MYSQL_DB_NAME)
	if err != nil {
		log.Fatalf("fail to create xorm Engine: %v", err)
	}
	err = orm.Sync(
		new(Category),
		new(Topic),
		new(Reply),
		new(Account),
	)
	//	orm.ShowSQL(true)
}

func GetOrm() *xorm.Engine {
	return orm
}
