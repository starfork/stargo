package broker

type Broker interface {
	Public(message Message) error
	Subscribe()
	UnSubscribe()
}

type Message struct {
	Header map[string]string
	Body   []byte
}
