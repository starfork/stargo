package testutil

import (
	"os"
	"time"
)

const (
	DefaultNATSURL  = "nats://localhost:4222"
	DefaultRedisURL = "localhost:6379"
	DefaultEtcdURL  = "localhost:2379"
	DefaultMySQLDSN = "root:stargo_test@tcp(localhost:3306)/stargo_test?charset=utf8mb4&parseTime=True"

	RetryInterval = 500 * time.Millisecond
	MaxRetries    = 20
)

func IsIntegrationTest() bool {
	return os.Getenv("INTEGRATION_TEST") == "1"
}

func GetEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func NATSURL() string {
	return GetEnvOrDefault("STARGO_TEST_NATS_URL", DefaultNATSURL)
}

func RedisAddr() string {
	return GetEnvOrDefault("STARGO_TEST_REDIS_ADDR", DefaultRedisURL)
}

func EtcdEndpoint() string {
	return GetEnvOrDefault("STARGO_TEST_ETCD_ENDPOINT", DefaultEtcdURL)
}

func MySQLDSN() string {
	return GetEnvOrDefault("STARGO_TEST_MYSQL_DSN", DefaultMySQLDSN)
}
