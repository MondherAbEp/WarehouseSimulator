package main

import "fmt"

func work(c constraints) {
	m, _ := createMatrix(c)

	if len(c.PalletTrucks) > 0 {
		pt := c.PalletTrucks[0]
		p, err := findClosestParcel(m, pt, c.Parcels)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(p)
		}
		t, err := findClosestTruck(m, pt, c.Trucks)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(t)
		}
	}
}
