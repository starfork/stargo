package mysql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/store"
	"github.com/starfork/stargo/util/ustring"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var TIME_LOCATION = "Asia/Shanghai" //上海
var TFORMAT = "2006-01-02T15:04:05+08:00"

//var log logger.Interface

type Mysql struct {
	db   *gorm.DB
	c    *config.StoreConfig
	conn *sql.DB
}

func NewMysql(config *config.StoreConfig) store.Store {
	return &Mysql{
		c: config,
	}
}

// Connect 初始化MySQLme
func (e *Mysql) Connect(confs ...*config.Config) {
	c := e.c
	if len(confs) > 0 {
		c = confs[0].Mysql
	}

	// if c.Timezome != "" {
	// 	TIME_LOCATION = c.Timezome
	// }
	// if c.Timeformat != "" {
	// 	TFORMAT = c.Timeformat
	// }
	c.User = ustring.Or(c.User, os.Getenv("MYSQL_USER"))
	c.Auth = ustring.Or(c.Auth, os.Getenv("MYSQL_PASSWD"))
	c.Host = ustring.Or(c.Host, os.Getenv("MYSQL_HOST"))
	c.Port = ustring.Or(c.Port, os.Getenv("MYSQL_PORT"))
	c.Name = ustring.Or(c.Name, os.Getenv("MYSQL_NAME"))

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Auth,
		c.Host,
		c.Port,
		c.Name,
	)

	conf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //全局采用单表名
			TablePrefix:   c.Prefix,
		},
	}
	if c.Debug {
		conf.Logger = logger.Default.LogMode(logger.Info)
	}
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), conf); err != nil {
		panic("Db Connect TO " + dsn + " With Error:" + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	//defer sqlDB.Close()
	sqlDB.SetMaxIdleConns(5)
	if c.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdle)
	}
	sqlDB.SetMaxOpenConns(10)
	if c.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpen)
	}
	//自己到项目里去注册
	//p := plugin.New(config)
	// if len(c.Plugins) > 0 {
	// 	RegisterPlugins(db, config, c.Plugins)
	// }
	e.db = db
	e.conn = sqlDB
}

func (e *Mysql) GetInstance(conf ...*config.Config) *gorm.DB {

	if len(conf) > 0 {
		e.Connect(conf...)
		return e.db
	}
	return e.db
}

func (e *Mysql) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *Mysql) Prefix(prefix string) string {

	e.db.NamingStrategy = schema.NamingStrategy{
		TablePrefix: prefix,
	}
	return prefix
}

// func (e *Mysql) RegisterPlugins(prefix string) string {

// 	e.db.NamingStrategy = schema.NamingStrategy{
// 		TablePrefix: prefix,
// 	}
// 	return prefix
// }
