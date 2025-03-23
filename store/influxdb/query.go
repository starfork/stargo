package influxdb

import (
	"fmt"
	"strings"
	"time"
)

type Query struct {
	bucket      string
	measurement string
	filters     []string
	drop        []string
	sortOrder   string
	limit       uint32
	offset      uint32
	tz          map[string]int64
	loc         string //时区
	pivot       string
	l           *time.Location
}

func NewQuery(bucket string, cloc ...string) *Query {
	loc := "Asia/Shanghai"

	if len(cloc) > 0 {
		loc = cloc[0]
	}
	l, _ := time.LoadLocation(loc) // 设置为中国上海时区（UTC+8）

	return &Query{
		bucket: bucket,
		l:      l,
		loc:    loc,
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
func (e *Query) Drop(tags []string) *Query {
	e.drop = tags
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
func (e *Query) Count() *Query {
	e.pivot = `group() |> count()`
	return e
}

func (e *Query) Build() string {
	query := fmt.Sprintf(`from(bucket: "%s")`, e.bucket)

	// 设置默认时间范围（最近 7 天）
	if e.tz == nil {
		now := time.Now().UTC().Unix()
		e.tz = map[string]int64{
			"from": now - 7*24*3600,
			"to":   now,
		}
	}

	// 时间范围转换为 UTC 格式
	startTime := time.Unix(e.tz["from"], 0).In(e.l).Format("2006-01-02T15:04:05Z")
	query += fmt.Sprintf(` |> range(start: %s`, startTime)
	if stop, exists := e.tz["to"]; exists {
		stopTime := time.Unix(stop, 0).In(e.l).Format("2006-01-02T15:04:05Z")
		query += fmt.Sprintf(`, stop: %s)`, stopTime)
	} else {
		query += `)`
	}

	// 过滤 Measurement
	if e.measurement != "" {
		query += fmt.Sprintf(` |> filter(fn: (r) => r["_measurement"] == "%s")`, e.measurement)
	}

	// 处理 Where 条件
	for _, filter := range e.filters {
		query += fmt.Sprintf(` |> filter(fn: (r) => %s)`, filter)
	}
	//不同条件可能需要去掉某些tag
	if len(e.drop) > 0 {
		query += fmt.Sprintf(`|> drop(columns: ["%s"])`, strings.Join(e.drop, `", "`))
	}

	// 统计时不排
	if !strings.Contains(e.pivot, "count()") {

		if e.sortOrder == "" {
			query += ` |> sort(columns: ["_time"], desc: true)`
		} else {
			query += fmt.Sprintf(` |> sort(columns: ["%s"]) |> sort(columns: ["_time"], desc: true)`, e.sortOrder)
		}

	}

	if e.pivot == "" {
		query += ` |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`
	} else {
		query += ` |>` + e.pivot
	}

	// 处理分页
	if e.limit > 0 {
		if e.offset > 0 {
			query += fmt.Sprintf(` |> limit(n: %d, offset: %d)`, e.limit, e.offset)
		} else {
			query += fmt.Sprintf(` |> limit(n: %d)`, e.limit)
		}
	}

	return query
}
