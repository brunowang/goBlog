package db

import (
	"log"
	"time"

	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	// 设置Mysql驱动名称
	_MYSQL_DRIVER = "mysql"
)

var (
	orm *xorm.Engine
	cfg *goconfig.ConfigFile
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
	cfg = loadConfig()
	_MYSQL_DB_SOURCE := getDbConf("uname") + ":" + getDbConf("pwd") + "@(" + getDbConf("host") + ":" + getDbConf("port") + ")/" + getDbConf("db_name")
	var err error
	orm, err = xorm.NewEngine(_MYSQL_DRIVER, _MYSQL_DB_SOURCE)
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

func loadConfig() *goconfig.ConfigFile {
	cfg, err := goconfig.LoadConfigFile("conf/db.conf")
	if err != nil {
		log.Fatalf("无法加载配置文件：%s", err)
	}
	return cfg
}

func getDbConf(key string) string {
	section := "mysql"
	value, err := cfg.GetValue(section, key)
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", key, err)
	}
	return value
}
