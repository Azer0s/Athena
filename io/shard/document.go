package shard

type Values map[string]interface{}
type Document struct {
	Id string
	Values
}

func (d *Document) Get(key string) interface{} {
	return d.Values[key]
}
