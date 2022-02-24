package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type warehouse struct {
	Width  int
	Height int
	Turns  int
}

type parcelType struct {
	name   string
	weight int
}

type parcel struct {
	X, Y int
	Type parcelType
}

type palletTruck struct {
	X, Y int
}

type truck struct {
	X, Y      int
	maxWeight int
	turns     int
}

type constraints struct {
	Warehouse    warehouse
	Parcels      map[string]parcel
	PalletTrucks map[string]palletTruck
	Trucks       map[string]truck
}

var errUnknownColor = errors.New("unknown Type")

func getParcelType(color string) (packageColor parcelType, err error) {
	color = strings.ToUpper(color)

	switch color {
	case "YELLOW":
		packageColor = parcelType{color, 100}
	case "GREEN":
		packageColor = parcelType{color, 200}
	case "BLUE":
		packageColor = parcelType{color, 500}
	default:
		err = fmt.Errorf("%w: %s", errUnknownColor, color)
	}
	return
}

func addPackage(constraints *constraints, values []string) {
	name := values[0]
	x, _ := strconv.Atoi(values[1])
	y, _ := strconv.Atoi(values[2])
	color, err := getParcelType(values[3])
	if err != nil {
		fmt.Print(err)
	}

	constraints.Parcels[name] = parcel{x, y, color}
}

func addPalletTruck(constraints *constraints, values []string) {
	name := values[0]
	x, _ := strconv.Atoi(values[1])
	y, _ := strconv.Atoi(values[2])

	constraints.PalletTrucks[name] = palletTruck{x, y}
}

func addTruck(constraints *constraints, values []string) {
	name := values[0]
	x, _ := strconv.Atoi(values[1])
	y, _ := strconv.Atoi(values[2])
	maxWeight, _ := strconv.Atoi(values[3])
	turns, _ := strconv.Atoi(values[4])

	constraints.Trucks[name] = truck{x, y, maxWeight, turns}
}

func assignWarehouse(constraints *constraints, line string) error {
	values := strings.Fields(line)

	if len(values) == 3 {
		width, err := strconv.Atoi(values[0])
		if err != nil || width < 1 {
			return fmt.Errorf("incorrect width: %s", values[0])
		}

		height, err := strconv.Atoi(values[1])
		if err != nil || height < 1 {
			return fmt.Errorf("incorrect height: %s", values[1])
		}

		turns, err := strconv.Atoi(values[2])
		if err != nil || turns < 1 {
			return fmt.Errorf("incorrect Turns: %s", values[2])
		}

		constraints.Warehouse = warehouse{width, height, turns}
	}
	return nil
}

func assignConstraint(constraints *constraints, line string) {
	values := strings.Fields(line)

	switch len(values) {
	case 4:
		addPackage(constraints, values)
	case 3:
		addPalletTruck(constraints, values)
	case 5:
		addTruck(constraints, values)
	}
}

func sanitizeLine(rawLine []byte) string {
	return strings.TrimSpace(string(rawLine))
}

func getConstraints(filename string) (constraints, error) {
	file, _ := os.Open(filename)
	defer file.Close()

	c := constraints{
		Warehouse:    warehouse{},
		Parcels:      make(map[string]parcel),
		PalletTrucks: make(map[string]palletTruck),
		Trucks:       make(map[string]truck),
	}

	reader := bufio.NewReader(file)

	rawLine, _, err := reader.ReadLine()
	line := sanitizeLine(rawLine)

	err = assignWarehouse(&c, line)
	if err != nil {
		return constraints{}, err
	}

	for {
		rawLine, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		line := sanitizeLine(rawLine)
		assignConstraint(&c, line)
	}

	return c, nil
}
