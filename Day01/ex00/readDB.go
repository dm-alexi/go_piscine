package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

func usage() {
	fmt.Println("Usage: readDB -f filename.xml | filename.json")
}

func main() {
	if len(os.Args) != 3 || os.Args[1] != "-f" {
		usage()
		return
	}
	var out []byte
	recipe, err := readUniversal(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		return
	}
	if strings.HasSuffix(os.Args[2], ".xml") {
		out, err = json.MarshalIndent(recipe, "", "    ")
	} else {
		out, err = xml.MarshalIndent(recipe, "", "    ")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		return
	}
	fmt.Println(string(out))
}
