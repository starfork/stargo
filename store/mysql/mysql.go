package mysql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/starfork/stargo/store"
	"github.com/starfork/stargo/util/ustring"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

//var log logger.Interface

type Mysql struct {
	db   *gorm.DB
	c    *store.Config
	conn *sql.DB
}

func NewMysql(config *store.Config) store.Store {
	return &Mysql{
		c: config,
	}
}

// Connect 初始化MySQLme
func (e *Mysql) connect(confs ...*store.Config) {
	c := e.c
	if len(confs) > 0 {
		c = confs[0]
	}

	var err error
	dsn := c.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local",
			ustring.Or(c.User, os.Getenv("MYSQL_USER")),
			ustring.Or(c.Auth, os.Getenv("MYSQL_PASSWD")),
			ustring.Or(c.Host, os.Getenv("MYSQL_HOST")),
			ustring.Or(c.Port, os.Getenv("MYSQL_PORT")),
			ustring.Or(c.Name, os.Getenv("MYSQL_NAME")),
		)
	}

	conf := &gorm.Config{

		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //全局采用单表名
			TablePrefix:   c.Prefix,
		},
	}
	if c.Debug {
		conf.Logger = logger.Default.LogMode(logger.Info)
	}
	if c.Tz1k != -1 {
		store.TZ1K = true
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
	if len(c.Plugins) > 0 {
		RegisterPlugins(db, c.Plugins)
	}
	e.db = db
	e.conn = sqlDB
}

func (e *Mysql) Instance(conf ...*store.Config) any {

	if len(conf) > 0 {
		e.connect(conf...)
		return e.db
	}
	if e.db == nil {
		e.connect()
	}
	return e.db
}

func (e *Mysql) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *Mysql) Prefix(prefix ...string) string {

	if len(prefix) > 0 {
		e.db.NamingStrategy = schema.NamingStrategy{
			TablePrefix: prefix[0],
		}
		return prefix[0]
	}

	return e.c.Prefix
}

// func (e *Mysql) RegisterPlugins(prefix string) string {

// 	e.db.NamingStrategy = schema.NamingStrategy{
// 		TablePrefix: prefix,
// 	}
// 	return prefix
// }
