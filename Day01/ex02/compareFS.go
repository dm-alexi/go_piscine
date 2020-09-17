package main

import (
	"bufio"
	"fmt"
	"os"
)

func usage() {
	fmt.Println("Usage: compareFS --old filename1 --new filename2")
}

func main() {
	if len(os.Args) != 5 || os.Args[1] != "--old" || os.Args[3] != "--new" {
		usage()
		return
	}
	m := make(map[string]bool)
	file, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m[scanner.Text()] = true
	}
	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	file.Close()
	file, err = os.Open(os.Args[4])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		if !m[scanner.Text()] {
			fmt.Println("ADDED", scanner.Text())
		} else {
			delete(m, scanner.Text())
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for key := range m {
		fmt.Println("REMOVED", key)
	}
}
