package path

import (
	"WarehouseSimulator/matrix"
	"errors"
	"fmt"
	"sort"
)

type Node struct {
	X      int
	Y      int
	cost   int
	parent *Node
}

type Path []*Node

type pathfinder struct {
	m      matrix.Matrix
	start  *Node
	end    *Node
	closed []*Node
	open   []*Node
	path   Path
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (p Path) String() (s string) {
	for i, n := range p {
		s += fmt.Sprintf("%d: X = %d, Y = %d\n", i, n.X, n.Y)
	}
	return
}

func (p *pathfinder) isStart(n *Node) bool {
	return n.X == p.start.X && n.Y == p.start.Y
}

func (p *pathfinder) isEnd(n *Node) bool {
	return n.X == p.end.X && n.Y == p.end.Y
}

func (p *pathfinder) isClosed(n *Node) bool {
	for _, c := range p.closed {
		if c.X == n.X && c.Y == n.Y {
			return true
		}
	}
	return false
}

func (p *pathfinder) removeOpen(n *Node) {
	for i, o := range p.open {
		if o.X == n.X && o.Y == n.Y {
			p.open = append(p.open[:i], p.open[i+1:]...)
		}
	}
}

func (p *pathfinder) addOpen(n *Node) {
	p.open = append(p.open, n)
}

func (p *pathfinder) addClosed(n *Node) {
	p.closed = append(p.closed, n)
}

func (p *pathfinder) addPath(n *Node) {
	p.path = append(p.path, n)
}

func (p *pathfinder) isOpen(n *Node) bool {
	for _, c := range p.open {
		if c.X == n.X && c.Y == n.Y {
			return true
		}
	}
	return false
}

func (p *pathfinder) isEmpty(n *Node) bool {
	if n.X < 0 || n.X >= p.m.Width || n.Y < 0 || n.Y >= p.m.Height {
		return false
	}

	if p.isClosed(n) || p.isOpen(n) {
		return false
	}

	c := p.m.Rows[n.Y][n.X]
	if c.Content != matrix.Empty {
		return false
	}

	return true
}

func (p *pathfinder) getLowestCostOpenNode() (*Node, error) {
	if len(p.open) == 0 {
		return &Node{}, errors.New("NO PATH")
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

func (p *pathfinder) computeCost(n *Node) int {
	return abs(p.end.X-n.X) + abs(p.end.Y-n.Y)
}

func (p *pathfinder) Find(parent *Node) {
	var n *Node

	// Top
	n = &Node{parent.X, parent.Y - 1, 0, parent}
	if p.isEmpty(n) {
		p.addOpen(n)
		p.open = append(p.open, n)
	}

	// Bottom
	n = &Node{parent.X, parent.Y + 1, 0, parent}
	if p.isEmpty(n) {
		p.addOpen(n)
	}

	// Left
	n = &Node{parent.X - 1, parent.Y, 0, parent}
	if p.isEmpty(n) {
		p.addOpen(n)
	}

	// Right
	n = &Node{parent.X + 1, parent.Y, 0, parent}
	if p.isEmpty(n) {
		p.addOpen(n)
	}

	parent, err := p.getLowestCostOpenNode()

	p.removeOpen(parent)
	p.addClosed(parent)

	if err != nil {
		fmt.Print(err)
	} else if p.isEnd(n) {
		p.end = n
		p.connectPath()
	} else {
		p.Find(parent)
	}
}

func Find(m matrix.Matrix, start *Node, end *Node) Path {
	p := &pathfinder{m, start, end, make([]*Node, 0), make([]*Node, 0), make([]*Node, 0)}

	p.closed = append(p.closed, start)
	p.Find(start)

	return p.path
}
