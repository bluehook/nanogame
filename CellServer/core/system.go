package core

func CreateSystem(t int, fn func(c *Chunck)) *System {
	return &System{
		Type:   t,
		Update: fn,
	}
}

type System struct {
	Type   int
	Update func(c *Chunck)
}
