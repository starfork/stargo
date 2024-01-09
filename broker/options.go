package broker

import (
	"context"
	"crypto/tls"
)

type Options struct {
	Context context.Context

	TLSConfig *tls.Config
	Addrs     []string
	Secure    bool
}

type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type SubscribeOptions struct {

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
}

type Option func(*Options)

type PublishOption func(*PublishOptions)

// PublishContext set context.
func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

type SubscribeOption func(*SubscribeOptions)

func NewOptions(opts ...Option) *Options {
	options := Options{
		Context: context.Background(),
		//Logger:  logger.DefaultLogger,
	}

	for _, o := range opts {
		o(&options)
	}

	return &options
}

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Addrs sets the host addresses to be used by the broker.
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// DisableAutoAck will disable auto acking of messages
// after they have been handled.
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

// Queue sets the name of the queue to share messages on.
func Queue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}

// Secure communication with the broker.
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// Specify TLS Config.
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// Logger sets the underline logger.
// func Logger(l logger.Logger) Option {
// 	return func(o *Options) {
// 		o.Logger = l
// 	}
// }

// SubscribeContext set context.
func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}
