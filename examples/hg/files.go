package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	// f := hg.NewHFile("somebigfile.txt")

	// // words := f.Iterator().Words()
	// // runes := f.Iterator().Runes()
	// // bytes := f.Iterator().Bytes()
	// lines := f.Iterator().Lines()

	// for lines.Next() {
	// 	fmt.Println(lines.HString())
	// }

	// if lines.Error() != nil {
	// 	fmt.Println(lines.Error())
	// }

	// if f.Error() != nil {
	// 	fmt.Println(f.Error())
	// }

	// or

	// f := hg.NewHFile("somebigfile.txt")

	// for line := f.Iterator().Lines(); line.Next(); {
	// 	fmt.Println(line.HString())
	// }

	// if f.Error() != nil {
	// 	fmt.Println(f.Error())
	// }

	//////////////////////////////////////////////////
	f := hg.NewHFile("some/dir/that/dont/exist/file.txt")
	defer f.Close()

	f.Append("one").Append("\n")
	f.Append("two").Append("\n")

	fmt.Printf("%s", f.Read())

	fmt.Printf("Name(): %v\n", f.Name())
	fmt.Printf("IsDir(): %v\n", f.IsDir())
	fmt.Printf("Size(): %v\n", f.Size())
	fmt.Printf("Mode(): %v\n", f.Mode())
	fmt.Printf("ModeTime(): %v\n", f.ModTime())

	fmt.Println(f.Exist())
	fmt.Println(f.HDir().Path())
	fmt.Println(f.Path())

	f.Rename("aaa/aaa/aaa/fff.txt").Copy(f.HDir().Join("copy_of_aaa.txt"))

	dir, file := f.Split()
	fmt.Println(dir.Path(), file.Path())

	fmt.Println(f.Ext())

	fmt.Println(f.MimeType())
}
