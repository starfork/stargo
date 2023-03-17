package request

import (
	"net"
	"net/http"
)

//获取IP。这个只根据remoteaddr获取的不是很准确
//要么nginx把IP转发到这个上面，要么继续完善各种可能从而获取准确的IP
func GetIp(req *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	return ip, err
}
