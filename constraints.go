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
	Name string
	X, Y int
	Type parcelType
}

type palletTruck struct {
	Name string
	X, Y int
}

type truck struct {
	Name      string
	X, Y      int
	maxWeight int
	turns     int
}

type constraints struct {
	Warehouse    warehouse
	Parcels      []parcel
	PalletTrucks []palletTruck
	Trucks       []truck
}

var errUnknownColor = errors.New("unknown Type")

func (c constraints) String() (s string) {
	s = "Constraints:\n"
	s += fmt.Sprintf("• Warehouse: width = %d, height = %d, Turns = %d\n", c.Warehouse.Width, c.Warehouse.Height, c.Warehouse.Turns)

	if len(c.Parcels) > 0 {
		s += "• Parcels:\n"
		for _, parcel := range c.Parcels {
			s += fmt.Sprintf("\t• %s: x = %d, y = %d, Type = %s\n", parcel.Name, parcel.X, parcel.Y, parcel.Type.name)
		}
	}
	if len(c.PalletTrucks) > 0 {
		s += "• PalletTrucks:\n"
		for _, palletTruck := range c.PalletTrucks {
			s += fmt.Sprintf("\t• %s: x = %d, y = %d\n", palletTruck.Name, palletTruck.X, palletTruck.Y)
		}
	}
	if len(c.Trucks) > 0 {
		s += "• Trucks:\n"
		for _, truck := range c.Trucks {
			s += fmt.Sprintf("\t• %s: x = %d, y = %d, maxWeight = %d, Turns = %d\n", truck.Name, truck.X, truck.Y, truck.maxWeight, truck.turns)
		}
	}
	return
}

func getParcelType(color string) (packageColor parcelType, err error) {
	color = strings.ToLower(color)

	switch color {
	case "yellow":
		packageColor = parcelType{color, 100}
	case "green":
		packageColor = parcelType{color, 200}
	case "blue":
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

	constraints.Parcels = append(constraints.Parcels, parcel{name, x, y, color})
}

func addPalletTruck(constraints *constraints, values []string) {
	name := values[0]
	x, _ := strconv.Atoi(values[1])
	y, _ := strconv.Atoi(values[2])

	constraints.PalletTrucks = append(constraints.PalletTrucks, palletTruck{name, x, y})
}

func addTruck(constraints *constraints, values []string) {
	name := values[0]
	x, _ := strconv.Atoi(values[1])
	y, _ := strconv.Atoi(values[2])
	maxWeight, _ := strconv.Atoi(values[3])
	turns, _ := strconv.Atoi(values[4])

	constraints.Trucks = append(constraints.Trucks, truck{name, x, y, maxWeight, turns})
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

func getConstraints(filename string) (constraints constraints, err error) {
	file, _ := os.Open(filename)
	defer file.Close()

	reader := bufio.NewReader(file)

	rawLine, _, err := reader.ReadLine()
	line := sanitizeLine(rawLine)

	err = assignWarehouse(&constraints, line)
	if err != nil {
		return
	}

	for {
		rawLine, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		line := sanitizeLine(rawLine)
		assignConstraint(&constraints, line)
	}

	return
}
