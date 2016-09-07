package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var rzScalar int64

func razor(f *os.File) (string, error) {
	fmt.Println("razoring")

	fstats, err := f.Stat()
	if err != nil {
		panic(err)
	}
	numBytes := fstats.Size()

	fi, err := f.Stat()
	if err != nil {
		return "nil", err
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return "nil", err
	}
	_, err = io.Copy(buf, f)
	if err != nil {
		return "nil", err
	}

	var bytesRazord int64
	bytesRazord = 0
	numLinesRazord := 0
	//remove lines until we have removed 1/rzScalar of file..
	for {
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return "nil", err
		}
		numLinesRazord++
		bytesRazord += int64(len(line))
		if bytesRazord >= numBytes/rzScalar {
			fmt.Printf("RESIZING COMPLETED!\nRemoved: %v bytes from file with size of %v\n", bytesRazord, numBytes)
			break
		}

	}
	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return "nil", err
	}
	nw, err := io.Copy(f, buf)
	if err != nil {
		return "nil", err
	}
	err = f.Truncate(nw)
	if err != nil {
		return "nil", err
	}
	err = f.Sync()
	if err != nil {
		return "nil", err
	}

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return "nil", err
	}

	return "completed", nil
}

func Razor(fptr *os.File) {
	rzScalar = 10
	str, err := razor(fptr)
	if err != nil {
		fmt.Println(str)
		panic(err)
	}
}
