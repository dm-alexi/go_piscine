package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type outFormat struct {
	showDir  bool
	showFile bool
	showLink bool
	ext      string
}

func printDir(path string, format *outFormat) {
	info, err := os.Lstat(path)
	if err != nil || !info.IsDir() {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	if format.showDir {
		fmt.Println(path)
	}
	l, _ := f.Readdirnames(-1)
	for _, name := range l {
		info, err := os.Lstat(path + name)
		if err != nil {
			continue
		}
		if info.IsDir() {
			printDir(path+name+"/", format)
		} else if info.Mode().IsRegular() && format.showFile {
			if strings.HasSuffix(name, format.ext) && len(name) > len(format.ext) {
				fmt.Println(path + name)
			}
		} else if info.Mode()&os.ModeSymlink != 0 && format.showLink {
			link, err := os.Readlink(path + name)
			if err == nil {
				_, err = os.Lstat(link)
			}
			if err != nil {
				link = "[broken]"
			}
			fmt.Println(path+name, "->", link)
		}
	}
}

func main() {
	var format outFormat

	flag.BoolVar(&format.showDir, "d", false, "show directories")
	flag.BoolVar(&format.showFile, "f", false, "show files")
	flag.BoolVar(&format.showLink, "sl", false, "show symbolic links")
	flag.StringVar(&format.ext, "ext", "", "restrict extensions (only with -f)")
	flag.Parse()
	if !format.showFile {
		format.ext = ""
	}
	if !(format.showDir || format.showFile || format.showLink) {
		format.showDir, format.showFile, format.showLink = true, true, true
	}
	if format.ext != "" {
		format.ext = "." + format.ext
	}
	paths := flag.Args()
	if len(paths) == 0 {
		paths = append(paths, "./")
	}
	for _, path := range paths {
		printDir(path, &format)
	}
}
