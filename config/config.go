package config

var (
	ENV_DEBUG      = "debug"
	ENV_PRODUCTION = "production"
)

// type Config struct {
// 	Deploy  string
// 	Monitor bool
// 	Base    *ServerConfig //如果各个服务么有单独设置，则公用
// 	Server  map[string]*ServerConfig
// }

// 公共配置模板
type Config struct {
	Environment string
	Org         string
	//ServerName string //服务名称--4-11改。通过app启动设置
	Port string //服务端口
	//ApiPort    string //api端口
	Timezome   string //时区设置
	Timeformat string

	Mysql      *MysqlConfig
	Redis      *RedisConfig
	MongoDb    *MongoDBConfig
	FileServer *FileServerConfig
	Log        *LogConfig
	Broker     *BrokerConfig
	Registry   *Registry
	Server     map[string]*Server //rpc server
}

// Mysql
type MysqlConfig struct {
	Name        string //数据库名字
	Host        string
	User        string
	Port        string
	Password    string
	Debug       bool //是否开启调试
	MaxIdle     int
	MaxOpen     int
	TablePrefix string
	Plugins     []string
}

// redis
type RedisConfig struct {
	Addr string //连接地址
	Auth string //认证
	Num  int    //库的数字
}

// Mongo
type MongoDBConfig struct {
	Host     string //地址
	Port     string //端口
	User     string //账户
	Password string //
	Monitor  bool   //监控
	DbName   string //库名
}

// 文件服务器配置
type FileServerConfig struct {
	PublicUrl  string //公共文件URL
	PrivateUrl string //私有文件URL
}

// log

type LogConfig struct {
	Target     string //日志输出目标。一般是console或者file
	LogFile    string //日志输出文件
	MaxSize    int    //日志文件最大尺寸
	MaxBackups int    //最大备份数
	MaxAge     int    //最大停留
	Level      int
}

type Registry struct {
	Environment string
	Org         string
	Name        string
	Addr        string //连接地址
	Auth        string //认证
	Num         int    //库的数字
}

type BrokerConfig struct {
	Name string
	Host string //连接地址
}

//Rpc Server

type Server struct {
	Name string
	Host string
	Port string
	Auth string //[keyfilepath]:[key]:
}
