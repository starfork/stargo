package geo

import (
	"math"
	"strconv"
	"strings"
)

//Point 点
type Point struct {
	lng float64 // 经度 longitude
	lat float64 // 纬度 latitude
}

// NewPoint returns a new Point populated by the passed in latitude (lat) and longitude (lng) values.
// 从原参考函数而来，调整了经纬度的顺序，符合mysql ST_Distance_Sphere 顺序。
func NewPoint(lng float64, lat float64) *Point {
	return &Point{lng: lng, lat: lat}
}

//NewTextPoint NewTextPoint 从逗号分割的点中创建point
func NewTextPoint(point string) *Point {
	p := strings.Split(point, ",")
	lng, _ := strconv.ParseFloat(strings.Trim(p[0], " "), 64)
	lat, _ := strconv.ParseFloat(strings.Trim(p[1], " "), 64)
	return &Point{lng: lng, lat: lat}
}

// Lat returns Point p's latitude.
func (p *Point) Lat() float64 {
	return p.lat
}

// Lng returns Point p's longitude.
func (p *Point) Lng() float64 {
	return p.lng
}

//Distance 点之间的距离
func Distance(p []*Point) float64 {
	if len(p) < 2 {
		return 0
	}
	var dist float64

	for i := 0; i < len(p)-1; i++ {
		dist = dist + distance(p[i], p[i+1])
	}
	return dist

}

//两点之间距离计算
//https://www.geodatasource.com/developers/go
//https://blog.csdn.net/juzipidemimi/article/details/104378053
func distance(p1, p2 *Point) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * p1.Lat() / 180)
	radlat2 := float64(PI * p2.Lat() / 180)

	theta := float64(p1.Lng() - p2.Lng())
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	dist = dist * (1.609344 * 1000) //本系统，单位都用米

	return dist
}
