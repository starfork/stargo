package store

type Config struct {
	Host string //地址
	Port string //端口
	User string //账户
	Name string //数据库名称，仓库名称，sqlite文件名等
	Auth string //认证/密码
	DSN  string //DSN连接

	Monitor bool //监控

	Plugins []string
	Debug   bool //是否开启调试
	MaxIdle int
	MaxOpen int
	Prefix  string //表前缀什么的
	Num     int    //连接标识数

}
