package main

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
)

type node struct {
	X      int
	Y      int
	F      float64
	G      float64
	H      float64
	parent *node
}

type path []*node

type pathfinder struct {
	m      matrix
	start  *node
	end    *node
	closed []*node
	open   []*node
}

type pathResult struct {
	i int
	p path
}

func (p path) String() (s string) {
	for i, n := range p {
		s += fmt.Sprintf("%d: X = %d, Y = %d\n", i, n.X, n.Y)
	}
	return
}

func (pf *pathfinder) isEnd(n *node) bool {
	return n.X == pf.end.X && n.Y == pf.end.Y
}

func (pf *pathfinder) isClosed(n *node) bool {
	for _, c := range pf.closed {
		if c.X == n.X && c.Y == n.Y {
			return true
		}
	}
	return false
}

func (pf *pathfinder) removeOpen(n *node) {
	for i, o := range pf.open {
		if o.X == n.X && o.Y == n.Y {
			pf.open[i] = pf.open[len(pf.open)-1]
			pf.open = pf.open[:len(pf.open)-1]
			break
		}
	}
}

func (pf *pathfinder) addOpen(n *node) {
	pf.open = append(pf.open, n)
}

func (pf *pathfinder) addClosed(n *node) {
	pf.closed = append(pf.closed, n)
}

func (pf *pathfinder) findOpen(n *node) *node {
	for _, c := range pf.open {
		if c.X == n.X && c.Y == n.Y {
			return c
		}
	}
	return nil
}

func (pf *pathfinder) isBorder(n *node) bool {
	return n.X < 0 || n.X >= pf.m.Width || n.Y < 0 || n.Y >= pf.m.Height
}

func (pf *pathfinder) isBlocked(n *node) bool {
	c := pf.m.Rows[n.Y][n.X]

	return c.Content == PalletTruck || c.Content == Parcel
}

func (pf *pathfinder) getLowestCostOpenNode() *node {
	sort.Slice(pf.open, func(i, j int) bool {
		return pf.open[i].F < pf.open[j].F
	})
	return pf.open[0]
}

func (pf *pathfinder) getPath(n *node) (p path) {
	for n.parent != nil {
		p = append([]*node{n}, p...)
		n = n.parent
	}
	return
}

func (pf *pathfinder) computeCost(n *node) {
	n.G = n.parent.G + 1
	n.H = math.Sqrt(float64((n.X - pf.end.X) ^ 2 + (n.Y - pf.end.Y) ^ 2))
	n.F = n.G + n.H
}

func (pf *pathfinder) find() (path, error) {
	pf.open = append(pf.open, pf.start)
	for len(pf.open) > 0 {
		parent := pf.getLowestCostOpenNode()
		pf.removeOpen(parent)

		nodes := [4]*node{
			{X: parent.X, Y: parent.Y - 1, parent: parent}, // TOP
			{X: parent.X - 1, Y: parent.Y, parent: parent}, // LEFT
			{X: parent.X, Y: parent.Y + 1, parent: parent}, // BOTTOM
			{X: parent.X + 1, Y: parent.Y, parent: parent}, // RIGHT
		}

		for _, n := range nodes {
			if !pf.isBorder(n) {
				if pf.isEnd(n) {
					return pf.getPath(n), nil
				}
				if !pf.isClosed(n) && !pf.isBlocked(n) {
					pf.computeCost(n)
					openNode := pf.findOpen(n)
					if openNode == nil || openNode.F > n.F {
						pf.addOpen(n)
					}
				}
			}
		}
		pf.addClosed(parent)
	}
	return nil, errors.New("no path found")
}

func findPath(m matrix, start node, end node) (path, error) {
	pf := pathfinder{m, &start, &end, make([]*node, 0), make([]*node, 0)}

	return pf.find()
}

func getAllPaths(m matrix, start node, ends []node) []pathResult {
	var paths []pathResult
	var wg sync.WaitGroup

	queue := make(chan pathResult, len(ends))

	for i, end := range ends {
		wg.Add(1)
		go func(i int, end node) {
			defer wg.Done()
			p, _ := findPath(m, start, end)
			queue <- pathResult{i, p}
		}(i, end)
	}

	wg.Wait()
	close(queue)

	for pResult := range queue {
		if len(pResult.p) > 0 {
			paths = append(paths, pResult)
		}
	}

	return paths
}

func findClosestParcel(m matrix, pt palletTruck, parcels []parcel) (parcel, error) {
	start := node{X: pt.X, Y: pt.Y}
	ends := make([]node, len(parcels))

	for i, p := range parcels {
		ends[i] = node{X: p.X, Y: p.Y}
	}

	paths := getAllPaths(m, start, ends)
	if len(paths) == 0 {
		return parcel{}, errors.New("no parcel found")
	}
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i].p) < len(paths[j].p)
	})
	return parcels[paths[0].i], nil
}

func findClosestTruck(m matrix, pt palletTruck, trucks []truck) (truck, error) {
	start := node{X: pt.X, Y: pt.Y}
	ends := make([]node, len(trucks))

	for i, t := range trucks {
		ends[i] = node{X: t.X, Y: t.Y}
	}

	paths := getAllPaths(m, start, ends)
	if len(paths) == 0 {
		return truck{}, errors.New("no truck found")
	}

	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i].p) < len(paths[j].p)
	})
	return trucks[paths[0].i], nil
}
