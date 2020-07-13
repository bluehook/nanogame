package core

import (
	"sync"
)

func CreateEntityManager() *EntityManager {
	return &EntityManager{
		bank: map[int]*Chunck{},
	}
}

func CreateSystemManager() *SystemManager {
	return &SystemManager{
		bank: map[int]*System{},
	}
}

//实体管理器
type EntityManager struct {
	bank map[int]*Chunck
}

//添加实体
func (mgr *EntityManager) AddEntity(e *Entity) {
	t := e.Type
	c, ok := mgr.bank[t]
	if ok {
		c.AddEntity(e)
	} else {
		c := CreateChunck(t)
		c.AddEntity(e)
		mgr.bank[t] = c
	}
}

func (mgr *EntityManager) RemoveEntity(e *Entity) {
	t := e.Type
	c, ok := mgr.bank[t]
	if ok {
		c.RemoveEntity(e)
	}
}

func (mgr *EntityManager) Find(t int) (*Chunck, bool) {
	c, ok := mgr.bank[t]
	if ok {
		return c, true
	} else {
		return nil, false
	}
}

func (mgr *EntityManager) Map(handler func(int, *Chunck)) {
	for k, c := range mgr.bank {
		handler(k, c)
	}
}

func (mgr *EntityManager) Len() int {
	return len(mgr.bank)
}

//系统管理器
type SystemManager struct {
	bank map[int]*System
}

func (mgr *SystemManager) Update(t int, c *Chunck, w *sync.WaitGroup) {
	defer w.Done()
	sys, ok := mgr.bank[t]
	if ok {
		sys.Update(c)
	}
}

func (mgr *SystemManager) AddSystem(sys *System) {
	mgr.bank[sys.Type] = sys
}

func (mgr *SystemManager) RemoveSystem(t int) {
	delete(mgr.bank, t)
}

func (mgr *SystemManager) Len() int {
	return len(mgr.bank)
}
