package command

import (
	"CellServer/core"
	"sync/atomic"
)

var roomId int64

func RoomID() int64 {
	return atomic.AddInt64(&roomId, 1)
}

//添加实体
type AddEntityCommand struct {
	Type int
	Data interface{}
}

func (c *AddEntityCommand) Execute() {
	entity := core.CreateEntity(c.Type, c.Data)
	core.GetWord().EntityMgr.AddEntity(entity)
}

//删除实体
type RemoveEntityCommand struct {
	Data interface{}
}

func (c RemoveEntityCommand) Execute() {
	core.GetWord().EntityMgr.RemoveEntity(c.Data.(*core.Entity))
}
