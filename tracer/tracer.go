package tracer

type Tracer interface {
	//SetUp(*Config) error
	Close() error
}
