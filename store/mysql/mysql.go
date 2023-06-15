package mysql

import (
	"database/sql"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/store/mysql/plugin"

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
	conn *sql.DB
}

// Connect 初始化MySQLme
func Connect(config *config.ServerConfig) *Mysql {
	c := config.Mysql
	if config.Timezome != "" {
		TIME_LOCATION = config.Timezome
	}

	var err error
	//dsn = "root:@tcp(127.0.0.1:3306)/zome_ucenter?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := c.User + ":" + c.Password + "@tcp(" +
		c.Host + ":" + c.Port + ")/" + c.Name +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	conf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //全局采用单表名
			TablePrefix:   c.TablePrefix,
		},
	}
	if c.Debug {
		conf.Logger = logger.Default.LogMode(logger.Info)
		//log = conf.Logger
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

	p := plugin.New(config)

	db.Callback().Query().After("gorm:find").Register("ossimage:after_query", p.AfterQuery)
	db.Callback().Update().Before("gorm:update").Register("ossimage:before_update", p.BeforeUpdate)

	return &Mysql{
		db:   db,
		conn: sqlDB,
	}
}

func (e *Mysql) GetInstance() *gorm.DB {

	return e.db
}

func (e *Mysql) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *Mysql) SetTablePrefix(prefix string) *gorm.DB {

	e.db.NamingStrategy = schema.NamingStrategy{
		TablePrefix: prefix,
	}
	return e.db
}
