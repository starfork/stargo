package mysql

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/starfork/stargo/util/config"
	"gorm.io/gorm"
)

type Store struct {
	table string
	db    *gorm.DB
}

func New(db *gorm.DB, table ...string) config.StoreInterface {
	s := &Store{
		db: db,
	}
	s.table = "config"
	if len(table) > 0 {
		s.table = table[0]
	}
	return s
}

// Setup  初始化数据库中的配置
func (e *Store) Load() []*config.KV {

	result := []*config.KV{}
	e.db.Table(e.table).Select("`key`", "`val`").Find(&result)
	return result

}

func (e *Store) Set(pfx string, value map[string]string) error {
	if len(value) == 0 {
		return errors.New("无效配置")
	}
	var buffer bytes.Buffer
	sql := "REPLACE INTO  `" + e.table + "` (`key`,`val`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for k, v := range value {
		k = strings.ToUpper(pfx) + "_" + strings.ToUpper(k)
		//cfg.Val[k] = v //important here
		buffer.WriteString(fmt.Sprintf("('%s','%s'),", k, v))
	}
	str := buffer.String()
	str = str[:len(str)-1]

	return e.db.Exec(str).Error

}
