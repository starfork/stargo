package influxdb

import (
	"fmt"
	"time"
)

type Query struct {
	bucket      string
	measurement string
	filters     []string
	sortOrder   string
	limit       uint32
	offset      uint32
	tz          map[string]int64
	loc         string //时区
	pivot       string
	l           *time.Location
}

func NewQuery(bucket string, loc ...string) *Query {
	l, _ := time.LoadLocation("Asia/Shanghai") // 设置为中国上海时区（UTC+8）

	if len(loc) > 0 {
		l, _ = time.LoadLocation(loc[0])
	}
	return &Query{
		bucket: bucket,
		l:      l,
	}
}

func (e *Query) Table(m string) *Query {
	e.measurement = m
	return e
}
func (e *Query) Tz(tz map[string]int64) *Query {
	e.tz = tz
	return e
}
func (e *Query) Pivot(tz string) *Query {
	e.pivot = tz
	return e
}

func (e *Query) Where(field, value string) *Query {
	if value != "" {
		e.filters = append(e.filters, fmt.Sprintf(`r["%s"] == "%s"`, field, value))
	}
	return e
}

func (e *Query) Order(order string) *Query {
	e.sortOrder = order
	return e
}

func (e *Query) Page(page, limit uint32) *Query {
	if page < 1 {
		page = 1
	}
	e.limit = limit
	e.offset = (page - 1) * limit
	return e
}

func (e *Query) Build() string {
	query := fmt.Sprintf(`from(bucket: "%s")`, e.bucket)

	// 设置默认时间范围为最近 30 天
	if e.tz == nil {
		now := time.Now().UTC().Unix()
		e.tz = map[string]int64{
			"from": now - 7*24*3600, // 默认7天。不建议太多
			"to":   now,             // 到当前时间
		}
	}

	// 添加 range 语句，确保时间为 UTC 格式
	startTime := time.Unix(e.tz["from"], 0).In(e.l).Format("2006-01-02T15:04:05Z")
	query += fmt.Sprintf(` |> range(start: %s`, startTime)
	if stop, exists := e.tz["to"]; exists {
		stopTime := time.Unix(stop, 0).In(e.l).Format("2006-01-02T15:04:05Z")
		query += fmt.Sprintf(`, stop: %s)`, stopTime)
	} else {
		query += `)`
	}

	// 添加 Measurement 过滤
	if e.measurement != "" {
		query += fmt.Sprintf(` |> filter(fn: (r) => r["_measurement"] == "%s")`, e.measurement)
	}

	// 添加字段过滤条件
	for _, filter := range e.filters {
		query += fmt.Sprintf(` |> filter(fn: (r) => %s)`, filter)
	}

	// 添加排序（如果设置了排序）
	if e.sortOrder != "" {
		query += fmt.Sprintf(` |> sort(columns: ["_time"], desc: %v)`, e.sortOrder == "desc")
	}
	if e.pivot == "" {
		query += ` |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`
	} else {
		query += ` |>` + e.pivot
	}
	// 添加分页（如果设置了分页）
	if e.offset > 0 {
		query += fmt.Sprintf(` |> offset(n: %d)`, e.offset)
	}
	if e.limit > 0 {
		query += fmt.Sprintf(` |> limit(n: %d)`, e.limit)
	}

	return query
}
