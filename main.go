package main

import (
	"os"
)

func main() {
	fptr, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fptr.Close()

	Razor(fptr)
}
