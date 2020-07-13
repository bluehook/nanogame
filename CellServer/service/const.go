package service

//实体类型
const (
	EntityTypeNetMessage   = iota + 1 //网络消息
	EntityTypeRoom                    //房间
	EntityTypeSessionJoin             //进入房间
	EntityTypeSessionLeave            //离开房间
	EntityTypeRoomPacket              //房间数据包
)
