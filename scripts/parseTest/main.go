package main

// file that reads from stdin and uses json.Unmarshal to parse the input to a rogue op

import (
	"fmt"
	"os"

	"github.com/heussd/pdftotext-go"
)

func main() {
	pdf, err := os.ReadFile("/Users/jpoz/Dropbox/Papers/Attention Is All You Need.pdf")
	if err != nil {
		panic(err)
	}

	pages, err := pdftotext.Extract(pdf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", pages)

	for _, page := range pages {
		fmt.Println(page.Content)
	}
}
