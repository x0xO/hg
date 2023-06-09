package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	for _, v := range hg.NewHDir("../*.go").Glob() {
		fmt.Println(v.Path())
	}

	fmt.Println("++++++++++++++")

	for _, v := range hg.NewHDir("..").ReadDir() {
		if v.IsDir() {
			continue
		}

		fmt.Println(v.Path())
	}

	d := hg.NewHDir("")
	fmt.Println(d.Path())
	fmt.Println(d.Exist())

	hg.NewHDir("./some/dir/that/dont/exist/").MkdirAll()

	d = hg.NewHDir("aaa").MkdirAll().Rename("bbb")
}
