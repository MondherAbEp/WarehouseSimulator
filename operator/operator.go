package operator

import (
	"WarehouseSimulator/constraints"
	"WarehouseSimulator/matrix"
	"WarehouseSimulator/path"
	"fmt"
)

func Work(c constraints.Constraints) {
	for turn := 0; turn < c.Warehouse.Turns; turn++ {
		m, _ := matrix.Create(c)
		pt := c.PalletTrucks[0]
		t := c.Trucks[0]

		p := path.Find(m, &path.Node{X: pt.X, Y: pt.Y}, &path.Node{X: t.X, Y: t.Y})
		fmt.Println(m)
		fmt.Println(p)
	}
}
