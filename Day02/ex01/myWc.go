package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
	"unicode/utf8"
)

func countLines(r *bufio.Reader) int {
	count := 0
	var line string
	var err error
	for err == nil {
		line, err = r.ReadString('\n')
		if len(line) > 0 {
			count++
		}
	}
	return count
}

func words(s string) int {
	if len(s) == 0 {
		return 0
	}
	count := 0
	if s[0] != ' ' && s[0] != '\t' && s[0] != '\n' {
		count++
	}
	for i, n := 1, len(s); i < n; i++ {
		if s[i] != ' ' && s[i] != '\t' && (s[i-1] == ' ' || s[i-1] == '\t') {
			count++
		}
	}
	return count
}

func countWords(r *bufio.Reader) int {
	count := 0
	var line string
	var err error
	for err == nil {
		line, err = r.ReadString('\n')
		count += words(line)
	}
	return count
}

func countCharacters(r *bufio.Reader) int {
	count := 0
	var line string
	var err error
	for err == nil {
		line, err = r.ReadString('\n')
		count += utf8.RuneCountInString(line)
	}
	return count
}

func getFunction() func(*bufio.Reader) int {
	var m, w, l bool
	flag.BoolVar(&m, "m", false, "characters")
	flag.BoolVar(&w, "w", false, "words")
	flag.BoolVar(&l, "l", false, "lines")
	flag.Parse()
	if (m && (w || l)) || (w && l) {
		return nil
	} else if m {
		return countCharacters
	} else if l {
		return countLines
	}
	return countWords
}

func main() {
	var wg sync.WaitGroup
	counterFunc := getFunction()
	if counterFunc == nil {
		fmt.Fprintln(os.Stderr, "Error: only one flag required")
		return
	}
	wg.Add(len(flag.Args()))
	for _, file := range flag.Args() {
		go func(s string, cf func(*bufio.Reader) int) {
			defer wg.Done()
			f, err := os.Open(s)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			defer f.Close()
			r := bufio.NewReader(f)
			fmt.Print(cf(r), "\t", s, "\n")
		}(file, counterFunc)
	}
	wg.Wait()
}
