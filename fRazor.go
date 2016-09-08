package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var rzScalar float64

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func razor(f *os.File) (*os.File, error) {
	fstats, err := f.Stat()
	checkErr(err)

	numBytes := fstats.Size()

	fi, err := f.Stat()
	checkErr(err)

	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, os.SEEK_SET)
	checkErr(err)

	_, err = io.Copy(buf, f)
	checkErr(err)

	var bytesRazord int64
	bytesRazord = 0
	numLinesRazord := 0
	//remove lines until we have removed ~1/rzScalar of file..
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			//exit on the last line
			fmt.Println("EOF")
			break
		}
		checkErr(err)

		fmt.Println("REMOVING: ", line)
		numLinesRazord++
		bytesRazord += int64(len(line))

		//we have removed the amount we need..so seek file and break from loop
		if float64(bytesRazord) >= float64(numBytes)*rzScalar {
			break
		}

	}

	_, err = f.Seek(0, os.SEEK_SET)
	checkErr(err)

	nw, err := io.Copy(f, buf)
	checkErr(err)

	err = f.Truncate(nw)
	checkErr(err)

	err = f.Sync()
	checkErr(err)

	//seek to end of file so our file pointer will be writing in the proper place
	_, err = f.Seek(fstats.Size(), os.SEEK_SET)
	checkErr(err)

	fmt.Printf("Resizing completed!\nRemoved: %v bytes from file with size of %v\n", bytesRazord, numBytes)
	return f, nil
}

//we return the fptr passed in so the outside program writing to the file is pointed to the proper, updated line in memory after line removal
func Razor(fptr *os.File, scalar float64) *os.File {
	rzScalar = scalar
	fptr, err := razor(fptr)
	if err != nil {
		panic(err)
	}

	return fptr
}
