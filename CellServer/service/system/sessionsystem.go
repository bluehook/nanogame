package system

import (
	"CellServer/core"
	"CellServer/message"
	"CellServer/net"
	"CellServer/service"
	"CellServer/service/entity"
	"log"
)

//进入房间
func SessionJoinSystem(c *core.Chunck) {
	if !c.Empty() {
		c.Map(func(e *core.Entity) {
			joinEntity := e.Data.(*entity.SessionJoinEntity)
			session := joinEntity.SessionData.(*net.Session)
			rid := session.String("RoomID")
			roomChunck, ok := core.GetWord().EntityMgr.Find(service.EntityTypeRoom)
			if ok {
				roomChunck.MapBreak(func(e *core.Entity) bool {
					roomEntity := e.Data.(*entity.RoomEntity)
					if roomEntity.ID == rid {
						//加入房间
						roomEntity.Sessions.Store(session.ID(), session)

						//发送房间ID列表
						packetIds := net.CreateWritePacket()
						packetIds.WriteInt8(message.MsgRoomIds)
						packetIds.WriteInt64(session.ID())
						packetIds.WriteInt64(session.ID())
						roomEntity.Sessions.Range(func(k, v interface{}) bool {
							client := v.(*net.Session)
							if client.ID() != session.ID() {
								packetIds.WriteInt64(client.ID())
							}
							return true
						})
						packetIds.WriteInt64(0)
						session.SendStream(packetIds.BytesSlice())

						//通知房间其他ID
						packet := net.CreateWritePacket()
						packet.WriteInt8(message.MsgRoomJoin)
						packet.WriteInt64(session.ID())
						roomEntity.Sessions.Range(func(k, v interface{}) bool {
							client := v.(*net.Session)
							if client.ID() != session.ID() {
								packet.WriteInt64(session.ID())
								client.SendStream(packet.BytesSlice())
							}
							return true
						})
						return true
					}
					return false
				})
			} else {
				session.KickSelf()
				log.Println("[", session.ID(), "] join room: '", rid, "' no find, KickSelf()")
			}
		})
		c.Clear()
	}
}

//离开房间
func SessionLeaveSystem(c *core.Chunck) {
	if !c.Empty() {
		c.Map(func(e *core.Entity) {
			leaveEntity := e.Data.(*entity.SessionLeaveEntity)
			session := leaveEntity.SessionData.(*net.Session)
			rid := session.String("RoomID")
			roomChunck, ok := core.GetWord().EntityMgr.Find(service.EntityTypeRoom)
			if ok {
				roomChunck.MapBreak(func(e *core.Entity) bool {
					roomEntity := e.Data.(*entity.RoomEntity)
					if roomEntity.ID == rid {
						//离开房间
						roomEntity.Sessions.Delete(session.ID())
						//通知房间其他玩家
						packet := net.CreateWritePacket()
						packet.WriteInt8(message.MsgRoomLeave)
						packet.WriteInt64(session.ID())
						roomEntity.Sessions.Range(func(k, client interface{}) bool {
							packet.WriteInt64(session.ID())
							client.(*net.Session).SendStream(packet.BytesSlice())
							return true
						})
						return true
					}
					return false
				})
			}
		})
		c.Clear()
	}
}
