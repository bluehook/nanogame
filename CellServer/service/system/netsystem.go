package system

import (
	"CellServer/core"
	"CellServer/message"
	"CellServer/net"
	"CellServer/service"
	"CellServer/service/command"
	"CellServer/service/entity"
	"log"
)

//网络消息处理
func NetMessageSystem(c *core.Chunck) {
	if !c.Empty() {
		c.Map(func(e *core.Entity) {
			defer func() {
				if err := recover(); err != nil {
					log.Println("NetMessageSystem Error: ", err)
				}
			}()

			//检查空包
			msgEntity := e.Data.(*entity.NetMessageEntity)
			omsg := msgEntity.Msg.([]byte)
			msg := net.CreateReadPacket(msgEntity.Msg.([]byte))
			header := msg.ReadInt8()
			if header != message.MsgPeer && header != message.MsgBroadcast {
				return
			} else {

				//检查单播,广播包
				var pType int8
				var not byte
				toId := msg.ReadInt64()
				session := msgEntity.Session.(*net.Session)
				roomid := session.String("RoomID")
				if header == message.MsgPeer {
					pType = 1
				} else {
					pType = 0
					not = msg.ReadByte()
				}

				//暂时不处理,添加延迟命令,让心跳更新处理
				packetEntity := entity.CreateRoomPacketEntity(pType, roomid, not != 0, session.ID(), toId, omsg)
				cmd := &command.AddEntityCommand{
					Type: service.EntityTypeRoomPacket,
					Data: packetEntity,
				}
				core.GetWord().Commit(cmd)
			}
		})
		c.Clear()
	}
}
