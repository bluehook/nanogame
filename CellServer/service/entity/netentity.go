package entity

func CreateNetMessageEntity(msg interface{}, s interface{}) *NetMessageEntity {
	return &NetMessageEntity{
		Msg:     msg,
		Session: s,
	}
}

//网络消息
type NetMessageEntity struct {
	Msg     interface{}
	Session interface{}
}
