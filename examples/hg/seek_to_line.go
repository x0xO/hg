package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg"
)

func main() {
	var (
		file       = hg.NewHFile("somebigfile.txt")
		position   int64
		content    hg.HString
		lineToRead = 10
	)

	position, content = file.SeekToLine(position, lineToRead)

	if file.Error() != nil {
		log.Fatal(file.Error())
	}

	fmt.Println(position)
	fmt.Println(content)
}
