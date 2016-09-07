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
	retStr := "hit error" //will reset to valid resizing str if we dont hit any errors...

	fstats, err := f.Stat()
	if err != nil {
		return retStr, err
	}
	numBytes := fstats.Size()

	fi, err := f.Stat()
	if err != nil {
		return retStr, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return retStr, err
	}
	_, err = io.Copy(buf, f)
	if err != nil {
		return retStr, err
	}

	var bytesRazord int64
	bytesRazord = 0
	numLinesRazord := 0
	//remove lines until we have removed ~1/rzScalar of file..
	for {
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return retStr, err
		}
		numLinesRazord++
		bytesRazord += int64(len(line))
		if bytesRazord >= numBytes/rzScalar {
			retStr = fmt.Sprintf("Resizing completed!\nRemoved: %v bytes from file with size of %v\n", bytesRazord, numBytes)
			return retStr, nil
		}

	}
	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return retStr, err
	}
	nw, err := io.Copy(f, buf)
	if err != nil {
		return retStr, err
	}
	err = f.Truncate(nw)
	if err != nil {
		return retStr, err
	}
	err = f.Sync()
	if err != nil {
		return retStr, err
	}

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return retStr, err
	}

	return retStr, nil
}

func Razor(fptr *os.File) {
	rzScalar = 10
	str, err := razor(fptr)
	if err != nil {
		fmt.Println(str)
		panic(err)
	}
}
