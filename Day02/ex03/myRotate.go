package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

func getNames(file, path string) (string, string) {
	tmp := strings.Split(file, "/")
	newname := path + tmp[len(tmp)-1]
	var archiveName string
	if n := strings.LastIndex(newname, "."); n > -1 {
		archiveName = newname[:n] + "_" + fmt.Sprint(time.Now().Unix()) + ".tag.gz"
	} else {
		archiveName = newname + "_" + fmt.Sprint(time.Now().Unix()) + ".tag.gz"
	}
	return newname, archiveName
}

func add(tw *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(tw, file)
	return err
}

func rotateLog(file, path string) error {
	newname, arcname := getNames(file, path)
	err := os.Rename(file, newname)
	if err != nil {
		return err
	}
	out, err := os.Create(arcname)
	if err != nil {
		return err
	}
	defer out.Close()
	gw := gzip.NewWriter(out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	err = add(tw, newname)
	if err != nil {
		return err
	}
	err = os.Remove(newname)
	return err
}

func main() {
	var dir string
	flag.StringVar(&dir, "a", ".", "set directory")
	flag.Parse()
	if dir[len(dir)-1] != '/' {
		dir = dir + "/"
	}
	var wg sync.WaitGroup
	wg.Add(len(flag.Args()))
	for _, file := range flag.Args() {
		go func(s string) {
			if err := rotateLog(s, dir); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			wg.Done()
		}(file)
	}
	wg.Wait()
}
