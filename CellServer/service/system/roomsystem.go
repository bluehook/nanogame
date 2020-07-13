package system

import (
	"CellServer/core"
	"CellServer/net"
	"CellServer/service"
	"CellServer/service/command"
	"CellServer/service/entity"
	"log"
)

//房间更新
func RoomSystem(c *core.Chunck) {
	if !c.Empty() {
		c.Map(func(e *core.Entity) {
			room := e.Data.(*entity.RoomEntity)
			empty := true
			room.Sessions.Range(func(k, v interface{}) bool {
				empty = false
				return false
			})
			if empty {
				room.EmptyTime += core.GetWord().HeartTime
			} else {
				room.EmptyTime = 0
			}
			if room.EmptyTime > int64(net.RoomEmptyKeepTime) {
				cmd := command.RemoveEntityCommand{
					Data: e,
				}
				core.GetWord().Commit(cmd)
				log.Println("Empty Room Timeout Remove: '", room.ID, "'")
			}
			//fmt.Println(room.ID, room.EmptyTime, int64(net.RoomEmptyKeepTime))
		})
	}
}

//客户端数据包处理
func RoomPacketSystem(c *core.Chunck) {
	if !c.Empty() {
		c.Map(func(e *core.Entity) {
			roomPacket := e.Data.(*entity.RoomPacketEntity)
			roomChunck, ok := core.GetWord().EntityMgr.Find(service.EntityTypeRoom)

			if ok {
				roomChunck.MapBreak(func(e *core.Entity) bool {
					room := e.Data.(*entity.RoomEntity)
					if room.ID == roomPacket.RoomID {
						if roomPacket.Type == 0 {
							//广播
							room.Sessions.Range(func(k, v interface{}) bool {
								client := v.(*net.Session)
								if roomPacket.NotSelf {
									client.SendStream(roomPacket.Data)
								} else if client.ID() != roomPacket.From {
									client.SendStream(roomPacket.Data)
								}
								return true
							})
						} else {
							//单播
							if client, ok := room.Sessions.Load(roomPacket.To); ok {
								client.(*net.Session).SendStream(roomPacket.Data)
							}
						}
						return true
					}
					return false
				})
			}

		})
		c.Clear()
	}
}
