package naming

type Config struct {
	Environment string

	Org  string
	Name string
	Host string //连接地址
	Auth string //认证
	Num  int    //库的数字
}
