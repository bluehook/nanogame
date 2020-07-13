package core

func CreateEntity(t int, data interface{}) *Entity {
	return &Entity{
		Type: t,
		Data: data,
	}
}

type Entity struct {
	Type int
	Id   int64
	Data interface{}
}
