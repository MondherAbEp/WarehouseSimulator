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

func (w *worker) checkDriver(c *constraints, drivers map[string]*driver) {
	p := c.Parcels[w.parcelName]
	d := drivers[w.truckName]

	if d.status == Gone {
		w.status = Waiting
	} else {
		fmt.Printf("%s LEAVE %s %s\n", w.name, w.parcelName, p.Type.name)
		d.totalWeight += p.Type.weight
		delete(c.Parcels, w.parcelName)
		if len(c.Parcels) == 0 {
			w.status = Blocked
		} else {
			w.assignParcel(c)
		}
	}
}

func (w *worker) move(c *constraints, drivers map[string]*driver) {
	if len(w.path) == 0 {
		switch w.status {
		case Waiting, TowardsTruck:
			w.checkDriver(c, drivers)
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

func (w *worker) work(c *constraints, drivers map[string]*driver) {
	switch w.status {
	case TowardsTruck, TowardsParcel:
		w.move(c, drivers)
	case Sleeping:
		w.assignParcel(c)
	case Waiting, Blocked:
		fmt.Printf("%s WAIT\n", w.name)
	}
}

func getWorkers(c constraints) map[string]*worker {
	workers := make(map[string]*worker, len(c.PalletTrucks))

	for name, pt := range c.PalletTrucks {
		w := &worker{status: Sleeping, name: name, pt: &pt}
		w.work(&c, nil)
		workers[name] = w
	}
	return workers
}
