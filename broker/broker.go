package broker

type Broker interface {
	Public()
	Subscribe()
	UnSubscribe()
}

type Message struct {
	Header map[string]string
	Body   []byte
}
