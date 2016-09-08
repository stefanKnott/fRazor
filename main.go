package main

import (
	"os"
)

func main() {
	fptr, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	//write to beginning of file - note that this should be removed via fRazor
	for {
		fptr.WriteString("logging\n")
		fptr2, stat := RazorCheck(fptr, 500, .2)
		if stat == "finished" {
			fptr = fptr2
			break
		}
	}
	defer fptr.Close()

	for {
		//log some more and extend maxBytes for file..
		//note that this should always be written at the END
		fptr.WriteString("MOAR LOGGING\n")
		fptr2, stat := RazorCheck(fptr, 700, .2)

		if stat == "finished" {
			fptr = fptr2
			break
		}
	}
}
