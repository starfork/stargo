package broker

import "time"

type Config struct {
	App  string
	Name string
	Host string //连接地址

	JetStream *JetStreamConfig
}

type JetStreamConfig struct {
	Enabled    bool
	StreamName string
	Subjects   []string
	MaxAge     time.Duration
	MaxBytes   int64
	MaxMsgs    int64
	Replicas   int
	MaxDeliver int
	AckWait    time.Duration
	DLQEnabled bool
	DLQName    string
}
