package main

import (
	"fmt"
	"os"
)

func start(c constraints) {
	workers := getWorkers(c)
	drivers := getDrivers(c)

	for turn := 1; turn <= c.Warehouse.Turns; turn++ {
		fmt.Printf("tour %d\n", turn)
		for _, w := range workers {
			w.work(&c, drivers)
		}
		for _, d := range drivers {
			d.work(c)
		}
		fmt.Println()
	}
	if len(c.Parcels) == 0 {
		fmt.Print("😎")
	} else {
		fmt.Print("🙂")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing file")
		return
	}

	c, err := getConstraints(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	start(c)
}
