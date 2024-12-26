package filemanager

// 云存储才需要
type Config struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	BucketName string `yaml:"bucket_name"`
}
