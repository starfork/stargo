package logger

type Config struct {
	Target     string //日志输出目标。一般是console或者file
	LogFile    string //日志输出文件
	MaxSize    int    //日志文件最大尺寸
	MaxBackups int    //最大备份数
	MaxAge     int    //最大停留
	Level      int
}
