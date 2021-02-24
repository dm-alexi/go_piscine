package main

import (
	"fmt"
)

func main() {
	presents := []Present{{5, 1}, {4, 5}, {3, 1}, {5, 2}}
	p, _ := getNCoolestPresents(presents, 3)
	fmt.Println(p)
}
