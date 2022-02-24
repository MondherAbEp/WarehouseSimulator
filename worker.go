package main

import "fmt"

const (
	TowardsParcel = iota
	TowardsTruck
	Waiting
	Sleeping
	Blocked
)

type worker struct {
	status     int
	name       string
	pt         *palletTruck
	parcelName string
	truckName  string
	path       path
}

func (w *worker) move(c *constraints) {
	if len(w.path) == 0 {
		switch w.status {
		case TowardsTruck:
			p := c.Parcels[w.parcelName]
			fmt.Printf("%s LEAVE %s %s\n", w.name, w.parcelName, p.Type.name)
			delete(c.Parcels, w.parcelName)
			w.assignParcel(c)
		case TowardsParcel:
			p := c.Parcels[w.parcelName]
			fmt.Printf("%s TAKE %s %s\n", w.name, w.parcelName, p.Type.name)
			w.assignTruck(c)
		}
	} else {
		node := w.path[0]
		w.path = w.path[1:]

		w.pt.X = node.X
		w.pt.Y = node.Y

		fmt.Printf("%s GO [%d,%d]\n", w.name, w.pt.X, w.pt.Y)
	}
}

func (w *worker) assignTruck(c *constraints) {
	m, _ := createMatrix(*c)

	name, path, err := findClosestTruck(m, *w.pt, c.Trucks)
	if err != nil {
		w.status = Blocked
	} else {
		w.truckName = name
		w.path = path
		w.status = TowardsTruck
	}
}

func (w *worker) assignParcel(c *constraints) {
	m, _ := createMatrix(*c)

	name, path, err := findClosestParcel(m, *w.pt, c.Parcels)
	if err != nil {
		w.status = Blocked
	} else {
		w.parcelName = name
		w.path = path
		w.status = TowardsParcel
	}
}

func (w *worker) work(c *constraints) {
	switch w.status {
	case TowardsTruck, TowardsParcel:
		w.move(c)
	case Sleeping:
		w.assignParcel(c)
	case Waiting:
		fmt.Printf("%s WAIT\n", w.name)
	}
}

func getWorkers(c constraints) map[string]*worker {
	workers := make(map[string]*worker, len(c.PalletTrucks))

	for name, pt := range c.PalletTrucks {
		w := &worker{status: Sleeping, name: name, pt: &pt}
		w.work(&c)
		workers[name] = w
	}
	return workers
}
