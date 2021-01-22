package main

import (
	"fmt"
)

func main() {
	fmt.Println("Test:")
	root := newNode(true)
	root.Left = newNode(true)
	root.Left.Left = newNode(true)
	root.Left.Right = newNode(false)
	root.Right = newNode(false)
	root.Right.Left = newNode(true)
	root.Right.Right = newNode(true)
	printTree(root)
	fmt.Println(unrollGarland(root))
}
