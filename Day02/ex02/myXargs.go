package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	s := strings.Join(strings.Fields(string(bytes)), " ")
	t := strings.Join(os.Args[1:], " ")
	cmd := exec.Command("bash", "-c", t+" "+s)
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error", err)
	}
	fmt.Print(string(out))
}
