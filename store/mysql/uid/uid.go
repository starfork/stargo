package uid

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Logger Log接口，如果设置了Logger，就使用Logger打印日志，如果没有设置，就使用内置库log打印日志
//var Logger logger

// ErrTimeOut 获取uid超时错误
var ErrTimeOut = errors.New("get uid timeout")

// UID struct
type UID struct {
	ch chan uint32 // id缓冲池

	min, max uint32 // id段最小值，最大值

	opt Options
}

// New 创建一个UID;len：缓冲池大小()
// db:数据库连接
// len：缓冲池大小(长度可控制缓存中剩下多少id时，去DB中加载)
func New(db *gorm.DB, opts ...Option) (*UID, error) {

	opt := DefaultOptions()
	for _, o := range opts {
		o(&opt)
	}
	opt.db = db
	lid := UID{
		ch:  make(chan uint32, opt.len),
		opt: opt,
	}
	go lid.productID()
	return &lid, nil
}

// Get 获取自增id,当发生超时，返回错误，避免大量请求阻塞，服务器崩溃
func (e *UID) Get() (uint32, error) {
	select {
	case <-time.After(1 * time.Second):
		return 0, ErrTimeOut
	case uid := <-e.ch:
		return uid, nil
	}
}

// productID 生产id，当ch达到最大容量时，这个方法会阻塞，直到ch中的id被消费
func (e *UID) productID() {
	e.reLoad()
	for {
		if e.min >= e.max {
			e.reLoad()
		}
		e.min++
		//过滤方法
		if len(e.opt.fun) > 0 {
			filter := e.opt.fun[0]
			if filter(e.min) != 0 {
				e.ch <- e.min
			}
		} else {
			e.ch <- e.min
		}

	}
}

// reLoad 在数据库获取id段，如果失败，会每隔一秒尝试一次
func (e *UID) reLoad() error {
	var err error
	for {
		err = e.getFromDB()
		if err == nil {
			return nil
		}

		// 数据库发生异常，等待一秒之后再次进行尝试
		if e.opt.logger != nil {
			e.opt.logger.Warnf("reload error %v", err)
		}
		time.Sleep(time.Second)
	}
}

// getFromDB 从数据库获取id段
func (e *UID) getFromDB() error {
	type result struct {
		MaxID uint32
		Step  uint32
	}
	var rs result

	tx := e.opt.db.Begin()
	defer tx.Rollback()

	if err := tx.Raw("SELECT max_id,step FROM "+e.opt.table+" WHERE business_id = ? FOR UPDATE", e.opt.id).Scan(&rs).Error; err != nil {
		return err
	}
	//步长过滤。避免productID多次调用db执行
	if len(e.opt.fun) > 1 {
		filter := e.opt.fun[1]
		rs.MaxID = filter(rs.MaxID, rs.Step)
	}

	if err := tx.Exec("UPDATE "+e.opt.table+" SET max_id = ? WHERE business_id = ?", rs.MaxID+rs.Step, e.opt.id).Error; err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	e.min = rs.MaxID
	e.max = rs.MaxID + rs.Step
	return nil
}
