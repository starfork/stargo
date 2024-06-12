package field

import (
	"context"
	"fmt"
	"strings"

	"github.com/twpayne/go-geom/encoding/wkb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Location  Create from customized data type
type Location struct {
	Point string
}

// GormDataType GormDataType
func (loc Location) GormDataType() string {
	return "geometry"
}

// Scan implements the sql.Scanner interface
func (loc *Location) Scan(v any) error {

	if v == nil {
		return nil
	}
	mysqlEncoding, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("did not scan: expected []byte but was %T", v)
	}
	var point wkb.Point
	point.Scan(mysqlEncoding[4:])
	co := point.Coords()
	loc.Point = fmt.Sprintf("POINT(%f %f)", co.X(), co.Y())

	return nil
}

// GormValue gorm value
func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	//return clause.Expr{}
	l := strings.Replace(loc.Point, ",", " ", -1)
	l = strings.Replace(l, "POINT(", "", -1)
	l = strings.Replace(l, ")", "", -1)
	return clause.Expr{
		SQL:  "ST_GeomFromText(?)",
		Vars: []any{"POINT(" + l + ")"},
	}
}
