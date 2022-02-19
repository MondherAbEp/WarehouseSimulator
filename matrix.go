package main

import (
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	Empty = iota
	Parcel
	PalletTruck
	Truck
)

type cell struct {
	Content int
	name    string
}

type row []cell

type matrix struct {
	Width  int
	Height int
	Rows   []row
}

func (m matrix) String() string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true

	for y := 0; y < m.Height; y++ {
		var row []interface{}
		for x := 0; x < m.Width; x++ {
			c := m.Rows[y][x]
			switch c.Content {
			case Empty:
				row = append(row, "  ")
			case Parcel:
				row = append(row, "📦")
			case PalletTruck:
				row = append(row, "👷")
			case Truck:
				row = append(row, "🚚")
			}
		}
		t.AppendRow(row)
	}
	return t.Render()
}

func allocateMatrix(m *matrix) {
	m.Rows = make([]row, m.Height)

	for y := 0; y < m.Height; y++ {
		m.Rows[y] = make(row, m.Width)
	}
}

func populateMatrix(m matrix, c constraints) {
	for _, parcel := range c.Parcels {
		m.Rows[parcel.Y][parcel.X] = cell{Parcel, parcel.Name}
	}
	for _, palletTruck := range c.PalletTrucks {
		m.Rows[palletTruck.Y][palletTruck.X] = cell{PalletTruck, palletTruck.Name}
	}
	for _, truck := range c.Trucks {
		m.Rows[truck.Y][truck.X] = cell{Truck, truck.Name}
	}
}

func createMatrix(c constraints) (m matrix, err error) {
	m.Width = c.Warehouse.Width
	m.Height = c.Warehouse.Height

	allocateMatrix(&m)
	populateMatrix(m, c)

	return
}
