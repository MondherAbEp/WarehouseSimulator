package matrix

import (
	"WarehouseSimulator/constraints"
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

type Matrix struct {
	Width  int
	Height int
	Rows   []row
}

func (m Matrix) String() string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true

	for y := 0; y < m.Height; y++ {
		var row []interface{}
		for x := 0; x < m.Width; x++ {
			cell := m.Rows[y][x]
			switch cell.Content {
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

func allocateMatrix(m *Matrix) {
	m.Rows = make([]row, m.Height)

	for y := 0; y < m.Height; y++ {
		m.Rows[y] = make(row, m.Width)
	}
}

func populateMatrix(m Matrix, c constraints.Constraints) {
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

func Create(c constraints.Constraints) (m Matrix, err error) {
	m.Width = c.Warehouse.Width
	m.Height = c.Warehouse.Height

	allocateMatrix(&m)
	populateMatrix(m, c)

	return
}
