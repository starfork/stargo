package sqlite

import (
	"database/sql"
	"os"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/util/ustring"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var TIME_LOCATION = "Asia/Shanghai" //上海
var TFORMAT = "2006-01-02T15:04:05+08:00"

//var log logger.Interface

type Sqlite struct {
	db   *gorm.DB
	conn *sql.DB
}

// Connect
func Connect(config *config.Config) *Sqlite {

	if config.Timezome != "" {
		TIME_LOCATION = config.Timezome
	}
	if config.Timeformat != "" {
		TFORMAT = config.Timeformat
	}
	c := config.Sqlite
	var err error

	name := ustring.Or(c.Name, os.Getenv("SQLITE_NAME"))

	conf := &gorm.Config{}
	if c.Debug {
		conf.Logger = logger.Default.LogMode(logger.Info)
	}
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(name), conf); err != nil {
		panic("db Connect TO " + name + " With Error:" + err.Error())
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

	//p := plugin.New(config)
	// if len(c.Plugins) > 0 {
	// 	RegisterPlugins(db, config, c.Plugins)
	// }

	return &Sqlite{
		db:   db,
		conn: sqlDB,
	}
}

func (e *Sqlite) GetInstance(conf ...*config.Config) *gorm.DB {

	if len(conf) > 0 {
		rs := Connect(conf[0])
		return rs.db
	}
	return e.db
}

func (e *Sqlite) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *Sqlite) SetTablePrefix(prefix string) *gorm.DB {

	e.db.NamingStrategy = schema.NamingStrategy{
		TablePrefix: prefix,
	}
	return e.db
}
