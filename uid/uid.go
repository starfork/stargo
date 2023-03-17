package uid

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type logger interface {
	Error(error)
}

// Logger Log接口，如果设置了Logger，就使用Logger打印日志，如果没有设置，就使用内置库log打印日志
var Logger logger

// ErrTimeOut 获取uid超时错误
var ErrTimeOut = errors.New("get uid timeout")

//UID struct
type UID struct {
	db         *gorm.DB    // 数据库连接
	businessID string      // 业务id
	ch         chan uint32 // id缓冲池

	min, max  uint32 // id段最小值，最大值
	CheckFunc []CheckFunc
}

type CheckFunc func(num ...uint32) uint32

// New 创建一个UID;len：缓冲池大小()
// db:数据库连接
// businessID：业务id
// len：缓冲池大小(长度可控制缓存中剩下多少id时，去DB中加载)
func New(db *gorm.DB, businessID string, len int) (*UID, error) {
	lid := UID{
		db:         db,
		businessID: businessID,
		ch:         make(chan uint32, len),
	}
	go lid.productID()
	return &lid, nil
}

// Get 获取自增id,当发生超时，返回错误，避免大量请求阻塞，服务器崩溃
func (u *UID) Get() (uint32, error) {
	select {
	case <-time.After(1 * time.Second):
		return 0, ErrTimeOut
	case uid := <-u.ch:
		return uid, nil
	}
}

// productID 生产id，当ch达到最大容量时，这个方法会阻塞，直到ch中的id被消费
func (u *UID) productID() {

	u.reLoad()

	for {
		if u.min >= u.max {
			u.reLoad()
		}
		u.min++

		//过滤方法
		if len(u.CheckFunc) > 0 {
			filter := u.CheckFunc[0]
			if filter(u.min) != 0 {
				u.ch <- u.min
			}
		} else {
			u.ch <- u.min
		}

	}
}

// reLoad 在数据库获取id段，如果失败，会每隔一秒尝试一次
func (u *UID) reLoad() error {
	var err error
	for {
		err = u.getFromDB()
		if err == nil {
			return nil
		}

		// 数据库发生异常，等待一秒之后再次进行尝试
		if Logger != nil {
			Logger.Error(err)
		} else {
			log.Println(err)
		}
		time.Sleep(time.Second)
	}
}

// getFromDB 从数据库获取id段
func (u *UID) getFromDB() error {
	type result struct {
		MaxID uint32
		Step  uint32
	}
	var rs result

	tx := u.db.Begin()
	defer tx.Rollback()

	err := tx.Raw("SELECT max_id,step FROM uid WHERE business_id = ? FOR UPDATE", u.businessID).Scan(&rs).Error
	// err := row.Scan(&maxID, &step)
	if err != nil {
		return err
	}
	//步长过滤。避免productID多次调用db执行
	if len(u.CheckFunc) > 1 {
		filter := u.CheckFunc[1]
		rs.MaxID = filter(rs.MaxID, rs.Step)
	}

	err = tx.Exec("UPDATE uid SET max_id = ? WHERE business_id = ?", rs.MaxID+rs.Step, u.businessID).Error
	if err != nil {
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	u.min = rs.MaxID
	u.max = rs.MaxID + rs.Step
	return nil
}
