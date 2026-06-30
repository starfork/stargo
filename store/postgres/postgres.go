package postgres

import (
	"database/sql"
	"time"

	"github.com/starfork/stargo/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func init() {
	store.Register("postgres", NewPostgres)
}

func NewPostgres(c *store.Config) store.Store {
	return &Postgres{c: c}
}

type Postgres struct {
	db   *gorm.DB
	c    *store.Config
	conn *sql.DB
}

func (e *Postgres) Instance(confs ...*store.Config) any {
	c := e.c
	if len(confs) > 0 {
		c = confs[0]
	}
	if err := e.connect(c); err != nil {
		return nil
	}
	return e.db
}

func (e *Postgres) InstanceE(confs ...*store.Config) (any, error) {
	c := e.c
	if len(confs) > 0 {
		c = confs[0]
	}
	if err := e.connect(c); err != nil {
		return nil, err
	}
	return e.db, nil
}

func (e *Postgres) GetInstance() *gorm.DB {
	return e.db
}

func (e *Postgres) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *Postgres) connect(c *store.Config) error {
	conf := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: c.Prefix,
		},
	}
	if !c.Debug {
		conf.Logger = gormlogger.Default.LogMode(gormlogger.Silent)
	}
	db, err := gorm.Open(postgres.Open(c.DSN), conf)
	if err != nil {
		return err
	}
	conn, err := db.DB()
	if err != nil {
		return err
	}
	conn.SetMaxIdleConns(5)
	if c.MaxIdle > 0 {
		conn.SetMaxIdleConns(c.MaxIdle)
	}
	conn.SetMaxOpenConns(10)
	if c.MaxOpen > 0 {
		conn.SetMaxOpenConns(c.MaxOpen)
	}
	conn.SetConnMaxLifetime(20 * time.Minute)
	e.db = db
	e.conn = conn
	return nil
}
