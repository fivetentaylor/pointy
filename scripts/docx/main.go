package main

import (
	"fmt"

	"code.sajari.com/docconv/v2"
)

func main() {
	res, err := docconv.ConvertPath("./scripts/docx/example.docx")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
