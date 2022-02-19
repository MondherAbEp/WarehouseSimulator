package main

import (
	"errors"
	"fmt"
	"sort"
)

type node struct {
	X      int
	Y      int
	cost   int
	parent *node
}

type path []*node

type pathfinder struct {
	m      matrix
	start  *node
	end    *node
	closed []*node
	open   []*node
	path   path
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (p path) String() (s string) {
	for i, n := range p {
		s += fmt.Sprintf("%d: X = %d, Y = %d\n", i, n.X, n.Y)
	}
	return
}

func (p *pathfinder) isEnd(n *node) bool {
	return n.X == p.end.X && n.Y == p.end.Y
}

func (p *pathfinder) isClosed(n *node) bool {
	for _, c := range p.closed {
		if c.X == n.X && c.Y == n.Y {
			return true
		}
	}
	return false
}

func (p *pathfinder) removeOpen(n *node) {
	for i, o := range p.open {
		if o.X == n.X && o.Y == n.Y {
			p.open = append(p.open[:i], p.open[i+1:]...)
		}
	}
}

func (p *pathfinder) addOpen(n *node) {
	p.open = append(p.open, n)
}

func (p *pathfinder) addClosed(n *node) {
	p.closed = append(p.closed, n)
}

func (p *pathfinder) addPath(n *node) {
	p.path = append(p.path, n)
}

func (p *pathfinder) isOpen(n *node) bool {
	for _, c := range p.open {
		if c.X == n.X && c.Y == n.Y {
			return true
		}
	}
	return false
}

func (p *pathfinder) isEmpty(n *node) bool {
	if n.X < 0 || n.X >= p.m.Width || n.Y < 0 || n.Y >= p.m.Height {
		return false
	}

	if p.isClosed(n) || p.isOpen(n) {
		return false
	}

	c := p.m.Rows[n.Y][n.X]
	if c.Content != Empty && c.Content != Truck {
		return false
	}

	return true
}

func (p *pathfinder) getLowestCostOpenNode() (*node, error) {
	if len(p.open) == 0 {
		return nil, errors.New("no path")
	}
	sort.Slice(p.open, func(i, j int) bool {
		return p.open[i].cost < p.open[j].cost
	})
	return p.open[0], nil
}

func (p *pathfinder) connectPath() {
	n := p.end

	for n.parent != nil {
		p.addPath(n)
		n = n.parent
	}

	for i, j := 0, len(p.path)-1; i < j; i, j = i+1, j-1 {
		p.path[i], p.path[j] = p.path[j], p.path[i]
	}
}

func (p *pathfinder) computeCost(n *node) int {
	return abs(p.end.X-n.X) + abs(p.end.Y-n.Y)
}

func (p *pathfinder) Find(parent *node) (err error) {
	nodes := [4]*node{
		{parent.X, parent.Y - 1, 0, parent}, // TOP
		{parent.X - 1, parent.Y, 0, parent}, // LEFT
		{parent.X, parent.Y + 1, 0, parent}, // BOTTOM
		{parent.X + 1, parent.Y, 0, parent}, // RIGHT
	}

	for _, n := range nodes {
		if p.isEmpty(n) {
			n.cost = p.computeCost(n)
			p.addOpen(n)
		}
	}

	parent, err = p.getLowestCostOpenNode()

	if err != nil {
		return
	}

	if p.isEnd(parent) {
		p.end = parent
		p.connectPath()
		return
	}
	p.removeOpen(parent)
	p.addClosed(parent)
	err = p.Find(parent)
	return
}

func findPath(m matrix, start *node, end *node) (path, error) {
	p := &pathfinder{m, start, end, make([]*node, 0), make([]*node, 0), make([]*node, 0)}

	p.closed = append(p.closed, start)
	err := p.Find(start)

	return p.path, err
}
