package main

import (
	"fmt"
)

func main() {
	empty := []Present{}
	fmt.Printf("Empy list. Capacity = %d, result = %v\n", 10, grabPresents(empty, 10))
	presents := []Present{{5, 1}, {4, 5}, {3, 1}, {5, 2}}
	for i := 0; i < 11; i++ {
		fmt.Printf("Capacity = %d, result = %v\n", i, grabPresents(presents, i))
	}
}
