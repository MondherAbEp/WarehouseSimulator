package main

import (
	"fmt"
	"os"
)

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

	work(c)
}
