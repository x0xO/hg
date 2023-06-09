package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	// f := hg.NewHFile("").TempFile("./", "*.txt").Write("some text")
	f := hg.NewHFile("").TempFile().Write("some text")
	fmt.Println(f.Path(), f.Read())

	fmt.Println(f.Read().Hash().MD5())

	f.Remove()
}
