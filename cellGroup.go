package balancer

import (
	"sync"
)

type CellGroup struct {
	mu     sync.Mutex
	node   Node
	cells  map[uint64]*cell
	load   uint64
	cRange Range
}

func NewCellGroup(n Node) CellGroup {
	return CellGroup{
		node:  n,
		cells: map[uint64]*cell{},
	}
}

func (cg *CellGroup) Node() Node {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.node
}

func (cg *CellGroup) SetNode(n Node) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.node = n
}

func (cg *CellGroup) SetRange(min, max uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.cRange = Range{
		Min: min,
		Max: max,
	}
}

// AddCell adds a cell to the cell group
// If autoremove flag is true, method calls CellGroup.RemoveCell of previous cell group.
// Flag is useful when CellGroup is altered and not refilled
func (cg *CellGroup) AddCell(c cell, autoremove bool) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	if cg == c.cg {
		return
	}
	cg.load += c.load
	cg.cells[c.id] = &c
	if c.cg != nil && autoremove {
		c.cg.RemoveCell(c.id)
	}
	c.cg = cg
}

// RemoveCell removes a cell from cell group.
func (cg *CellGroup) RemoveCell(id uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	delete(cg.cells, id)
	return
}

func (cg *CellGroup) TotalLoad() uint64 {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.load
}

func (cg *CellGroup) addLoad(l uint64) {
	cg.load += l
}
