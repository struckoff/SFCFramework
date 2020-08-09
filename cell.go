package balancer

import (
	"sync"
)

type cell struct {
	id   uint64
	mu   sync.Mutex
	load uint64
	off  map[string]uint64 // location of Relocated DataItem. DataItem.ID -> cell.ID
	dis  map[string]uint64 // data items in cell
	cg   *CellGroup
}

func NewCell(id uint64, cg *CellGroup, load uint64) *cell {
	c := cell{
		id:   id,
		load: load,
		cg:   cg,
		off:  make(map[string]uint64),
		dis:  make(map[string]uint64),
	}
	//found := false
	//for i := range cgs {
	//	if id >= cgs[i].cRange.Min && id < cgs[i].cRange.Max {
	//		found = true
	//		c.cg = &cgs[i]
	//		break
	//	}
	//}
	//if !found { //? May be this could be c.cg == nil
	//	return nil, errors.New("unable to bind cell to cell group")
	//}
	return &c
}

func (c *cell) ID() uint64 {
	return c.id
}

func (c *cell) SetGroup(cg *CellGroup) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cg = cg
}

func (c *cell) Group() *CellGroup {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cg
}

func (c *cell) Load() (load uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, l := range c.dis {
		load += l
	}
	return load
}

func (c *cell) Truncate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dis = make(map[string]uint64)
	c.load = 0
}

func (c *cell) Add(d DataItem) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.load += d.Size()
	c.dis[d.ID()] = d.Size()
	c.cg.AddLoad(d.Size())
	return nil
}

func (c *cell) Remove(d DataItem) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.off[d.ID()]; ok {
		delete(c.off, d.ID())
	}
	if _, ok := c.dis[d.ID()]; ok {
		c.load -= d.Size()
		delete(c.dis, d.ID())
		c.cg.RemoveLoad(d.Size())
	}
	return nil
}

func (c *cell) Relocate(d DataItem, ncID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.off[d.ID()] = ncID
	if _, ok := c.dis[d.ID()]; ok {
		delete(c.dis, d.ID())
		//c.cg.removeLoad(d.Size())
		c.load -= d.Size()
	}
}

func (c *cell) Relocated(did string) (uint64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cid, ok := c.off[did]
	return cid, ok
}
