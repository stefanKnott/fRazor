package main

import (
	"os"
)

func main() {
	fptr, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0666)
	//write to beginning of file - note that this should be removed via fRazor
	fptr.WriteString("OMG THIS FILE HAS BEEN OPENED WHO DID THIS?")
	if err != nil {
		panic(err)
	}
	defer fptr.Close()

	fptr = Razor(fptr, .5)
	//write to end -- the fptr will be directed at EOF upon finishing Razor
	fptr.WriteString("\nIF THIS ISNT AT THE END OF THE FILE IM GOING TO FREAK")
}
