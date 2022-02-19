package operator

import (
	"WarehouseSimulator/constraints"
	"WarehouseSimulator/matrix"
	"WarehouseSimulator/path"
	"fmt"
	"log"
	"time"
)

func Work(c constraints.Constraints) {
	m, _ := matrix.Create(c)
	pt := &c.PalletTrucks[0]
	t := c.Trucks[0]

	start := time.Now()
	p, err := path.Find(m, &path.Node{X: pt.X, Y: pt.Y}, &path.Node{X: t.X, Y: t.Y})
	elapsed := time.Since(start)
	log.Printf("Pathfinder took %s", elapsed)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(p)
	}
}
