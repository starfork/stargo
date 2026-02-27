package naming

type Config struct {
	Org         string
	Environment string
	Scheme      string //那种驱动类型，redis，etcd

	Host string //连接地址. 多个用逗号隔开
	Auth string //认证。多个用逗号隔开。于host一一对应
	Num  int    //库的数字
	Ttl  int64  //过期时间
}
