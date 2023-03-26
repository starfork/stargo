package jwt

type Data map[string]interface{}

func NewData() Data {
	d := make(Data)
	return d
}

func (d Data) Set(name string, value interface{}) Data {

	d[name] = value
	return d
}

func (d Data) Uid() uint32 {

	return d["Uid"].(uint32)
}
