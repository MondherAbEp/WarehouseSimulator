package main

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

func allocateMatrix(m *matrix) {
	m.Rows = make([]row, m.Height)

	for y := 0; y < m.Height; y++ {
		m.Rows[y] = make(row, m.Width)
	}
}

func populateMatrix(m matrix, c constraints) {
	for name, parcel := range c.Parcels {
		m.Rows[parcel.Y][parcel.X] = cell{Parcel, name}
	}
	for name, palletTruck := range c.PalletTrucks {
		m.Rows[palletTruck.Y][palletTruck.X] = cell{PalletTruck, name}
	}
	for name, truck := range c.Trucks {
		m.Rows[truck.Y][truck.X] = cell{Truck, name}
	}
}

func createMatrix(c constraints) (m matrix, err error) {
	m.Width = c.Warehouse.Width
	m.Height = c.Warehouse.Height

	allocateMatrix(&m)
	populateMatrix(m, c)

	return
}
