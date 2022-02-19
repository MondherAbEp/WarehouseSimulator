package main

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

type parcelPathResult struct {
	i int
	p path
}

func findPaths(m matrix, pt palletTruck, parcels []parcel) []parcelPathResult {
	var paths []parcelPathResult
	var wg sync.WaitGroup

	queue := make(chan parcelPathResult, len(parcels))

	for i, p := range parcels {
		wg.Add(1)

		go func(i int, p parcel) {
			defer wg.Done()
			parcelPath, _ := findPath(m, &node{X: pt.X, Y: pt.Y}, &node{X: p.X, Y: p.Y})
			queue <- parcelPathResult{i, parcelPath}
		}(i, p)
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
	paths := findPaths(m, pt, parcels)

	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i].p) < len(paths[j].p)
	})

	if len(paths) == 0 {
		return parcel{}, errors.New("no parcel found")
	}
	return parcels[paths[0].i], nil
}

func work(c constraints) {
	m, _ := createMatrix(c)
	pt := c.PalletTrucks[0]

	start := time.Now()
	p, err := findClosestParcel(m, pt, c.Parcels)
	elapsed := time.Since(start)
	fmt.Printf("findClosestParcel took %s\n", elapsed)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(p)
	}
}
