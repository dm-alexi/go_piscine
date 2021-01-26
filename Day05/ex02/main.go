package main

import (
	"fmt"
)

func main() {
	presents := []Present{{5, 1}, {4, 5}, {3, 1}, {5, 2}}
	fmt.Println(getCoolestPresent(presents))
}
