package naming

type Config struct {
	Environment string
	Scheme      string //那种驱动类型，redis，etcd

	Org  string
	Host string //连接地址
	Auth string //认证
	Num  int    //库的数字
	Ttl  int64  //过期时间
}
