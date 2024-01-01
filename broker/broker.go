package broker

type Broker interface {
	Public()
	Subscribe()
}

type Message struct {
	Header map[string]string
	Body   []byte
}
