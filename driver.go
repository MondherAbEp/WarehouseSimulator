package main

import "fmt"

const (
	Parked = iota
	Gone
)

type driver struct {
	status         int
	name           string
	t              *truck
	totalWeight    int
	remainingTurns int
}

func (d *driver) checkWeight(c constraints) {
	if (float64(d.totalWeight)/float64(d.t.maxWeight) > 0.5) ||
		(d.totalWeight > 0 && len(c.Parcels) == 0) {
		d.remainingTurns = d.t.turns
		d.status = Gone
		fmt.Printf("%s GONE %d/%d\n", d.name, d.totalWeight, d.t.maxWeight)
	} else {
		fmt.Printf("%s WAITING %d/%d\n", d.name, d.totalWeight, d.t.maxWeight)
	}
}

func (d *driver) work(c constraints) {
	switch d.status {
	case Gone:
		if d.remainingTurns == 0 {
			d.totalWeight = 0
			d.status = Parked
			fmt.Printf("%s WAITING %d/%d\n", d.name, d.totalWeight, d.t.maxWeight)
		} else {
			fmt.Printf("%s GONE %d/%d\n", d.name, d.totalWeight, d.t.maxWeight)
			d.remainingTurns--
		}
	case Parked:
		d.checkWeight(c)
	}
}

func getDrivers(c constraints) map[string]*driver {
	drivers := make(map[string]*driver, len(c.Trucks))

	for name, t := range c.Trucks {
		w := &driver{status: Parked, name: name, t: &t}
		drivers[name] = w
	}
	return drivers
}
