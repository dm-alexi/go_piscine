package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func usage() {
	fmt.Print("Usage: d00 [arguments]\n\t-h\thelp\n\tmean\tdisplay mean\n\tmedian\tdisplay median\n\tmode\tdisplay mode\n\tsd\tdisplay standard deviation\n")
}

func getMedian(arr []int) float64 {
	n := len(arr)
	if n%2 == 1 {
		return float64(arr[n/2])
	}
	return (float64(arr[n/2]) + float64(arr[n/2-1])) / 2
}

func getMean(arr []int) float64 {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return float64(sum) / float64(len(arr))
}

func getMode(arr []int) int {
	mode, count := arr[0], 1
	for i, c, n := 1, 1, len(arr); i < n; i++ {
		if arr[i] == arr[i-1] {
			c++
			if c > count {
				mode, count = arr[i], c
			}
		} else {
			c = 1
		}
	}
	return mode
}

func main() {
	format := 0
	for _, v := range os.Args[1:] {
		switch v {
		case "-h":
			usage()
			return
		case "mean":
			format |= 1
		case "median":
			format |= 2
		case "mode":
			format |= 4
		case "sd":
			format |= 8
		default:
			fmt.Fprintln(os.Stderr, "Error: invalid argument:", v)
			return
		}
	}
	if format == 0 {
		format = 15
	}
	var arr []int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if num, err := strconv.Atoi(scanner.Text()); err == nil && num >= -100000 && num <= 100000 {
			arr = append(arr, num)
		} else {
			fmt.Fprintln(os.Stderr, "Invalid input:", scanner.Text())
			return
		}
	}
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", scanner.Err())
	}
	if len(arr) == 0 {
		usage()
		return
	}
	sort.Ints(arr)
	mean := getMean(arr)
	fmt.Println("```")
	if format&1 > 0 {
		fmt.Printf("Mean: %.2f\n", mean)
	}
	if format&2 > 0 {
		fmt.Printf("Median: %.2f\n", getMedian(arr))
	}
	if format&4 > 0 {
		fmt.Printf("Mode: %d\n", getMode(arr))
	}
	if format&8 > 0 {
		sd := 0.0
		for _, v := range arr {
			sd += (float64(v) - mean) * (float64(v) - mean)
		}
		sd = math.Sqrt(sd / float64(len(arr)))
		fmt.Printf("SD: %.2f\n", sd)
	}
	fmt.Println("```")
}
