package main

import (
	"fmt"
	"log"
	"time"
)

func work(c constraints) {
	m, _ := createMatrix(c)
	pt := &c.PalletTrucks[0]
	t := c.Trucks[0]

	start := time.Now()
	p, err := findPath(m, &node{X: pt.X, Y: pt.Y}, &node{X: t.X, Y: t.Y})
	elapsed := time.Since(start)
	log.Printf("Pathfinder took %s", elapsed)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(p)
	}
}
