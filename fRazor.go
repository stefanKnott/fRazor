package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func razor(f *os.File, fileSize int64, rzScalar float64) (*os.File, error) {
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

		//we have removed the amount we need
		if float64(bytesRazord) >= float64(fileSize)*rzScalar {
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
	_, err = f.Seek(fileSize, os.SEEK_SET)
	checkErr(err)

	fmt.Printf("Resizing completed!\nRemoved: %v bytes from file with size of %v\n", bytesRazord, fileSize)
	return f, nil
}

//we return the fptr passed in so the outside program writing to the file is pointed to the proper EOF spot after line removals
func RazorCheck(fptr *os.File, maxBytes int64, scalar float64) (*os.File, string) {
	fstats, err := fptr.Stat()
	checkErr(err)

	//check if file is at or over the size limit
	if fstats.Size() >= maxBytes {
		fptr, err = razor(fptr, fstats.Size(), scalar)
		checkErr(err)
		return fptr, "finished"
	}

	return fptr, "not finished"
}
