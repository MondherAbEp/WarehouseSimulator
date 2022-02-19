package main

import (
	"WarehouseSimulator/constraints"
	"WarehouseSimulator/operator"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing file")
		return
	}

	c, err := constraints.Get(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	operator.Work(c)
}
