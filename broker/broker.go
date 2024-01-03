package broker

type Broker interface {
	Public(Message) error
	Subscribe()
	UnSubscribe()
}

type Message struct {
	Header map[string]string
	Body   []byte
}
