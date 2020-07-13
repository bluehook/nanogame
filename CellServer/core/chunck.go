package core

func CreateChunck(t int) *Chunck {
	c := &Chunck{
		Type: t,
		data: map[int64]*Entity{},
	}
	return c
}

//实体块
type Chunck struct {
	Type    int
	data    map[int64]*Entity
	idCount int64
}

func (c *Chunck) AddEntity(e *Entity) {
	e.Id = c.newId()
	c.data[e.Id] = e
}

func (c *Chunck) RemoveEntity(e *Entity) {
	delete(c.data, e.Id)
}

func (c *Chunck) Find(id int64) (*Entity, bool) {
	e, ok := c.data[id]
	if ok {
		return e, true
	} else {
		return nil, false
	}
}

func (c *Chunck) Clear() {
	c.data = map[int64]*Entity{}
}

func (c *Chunck) Empty() bool {
	return len(c.data) == 0
}

func (c *Chunck) Map(handler func(entity *Entity)) {
	for _, e := range c.data {
		handler(e)
	}
}

func (c *Chunck) MapBreak(handler func(entity *Entity) bool) {
	for _, e := range c.data {
		if handler(e) {
			return
		}
	}
}

func (c *Chunck) Len() int {
	return len(c.data)
}

func (c *Chunck) newId() int64 {
	c.idCount += 1
	return c.idCount
}
