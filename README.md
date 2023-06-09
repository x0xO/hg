<p align="center">
  <img src="https://user-images.githubusercontent.com/65846651/229838021-741ff719-8c99-45f6-88d2-1a32927bd863.png">
</p>

# 🤪 Humango: Go Crazy, Go Human, Go Nuts!
[![Go Reference](https://pkg.go.dev/badge/github.com/x0xO/hg.svg)](https://pkg.go.dev/github.com/x0xO/hg)
[![Go Report Card](https://goreportcard.com/badge/github.com/x0xO/hg)](https://goreportcard.com/report/github.com/x0xO/hg)

Introducing Humango, the wackiest Go package on the planet, created to make your coding experience an absolute riot! With Humango, you can forget about dull and monotonous code; we're all about turning the mundane into the insanely hilarious. It's not just a bicycle; it's almost a motorcycle 🤣!

## 🎉 What's in the box?
1. 📖 **Readable syntax**: Boring code is so yesterday! Humango turns your code into a party by blending seamlessly with Go and making it super clean and laughably maintainable.
2. 🔀 **Encoding and decoding:** Juggling data formats? No problemo! Humango's got your back with __Base64__, __URL__, __Gzip__, and __Rot13__ support. Encode and decode like a pro!
3. 🔒 **Hashing extravaganza:** Safety first, right? Hash your data with __MD5__, __SHA1__, __SHA256__, or __SHA512__, and enjoy peace of mind while Humango watches over your bytes.
4. 📁 **File and directory shenanigans:** Create, read, write, and dance through files and directories with Humango's fun-tastic functions. Trust us, managing files has never been this entertaining.
5. 🌈 **Data type compatibility:** Strings, integers, floats, bytes, slices, maps, you name it! Humango is the life of the party, mingling with all your favorite data types.
6. 🔧 **Customize and extend:** Need something extra? Humango is your best buddy, ready to be extended or modified to suit any project.
7. 📚 **Docs & examples:** We're not just about fun and games, we've got detailed documentation and examples that'll have you smiling from ear to ear as you learn the Humango way.

Take your Go projects to a whole new level of excitement with Humango! It's time to stop coding like it's a chore and start coding like it's a celebration! 🥳

# Examples

Generate a securely random string.

<table>
<tr>
<th><code>stdlib</code></th>
<th><code>hg</code></th>
</tr>
<tr>
<td>

```go
func main() {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	length := 10

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return
	}

	for i, v := range b {
		b[i] = charset[v%byte(len(charset))]
	}

	result := string(b)
	fmt.Println(result)
}
```
</td>
<td>

```go
func main() {
	result := hg.NewHString().Random(10)
	fmt.Println(result)
}
```
</td>
</tr>
</table>

GetOrDefault returns the value for a key. If the key does not exist, returns the default value
instead. This function is useful when you want to provide a fallback value for keys that may not
be present in the map.

<table>
<tr>
<th><code>stdlib</code></th>
<th><code>hg</code></th>
</tr>
<tr>
<td>

```go
func main() {
	md := map[int][]int{}

	for i := 0; i < 5; i++ {
		value, ok := md[i]
		if !ok {
			value = []int{}
		}

		md[i] = append(value, i)
	}

	fmt.Println(md)
}
```
</td>
<td>

```go
func main() {
	md := hg.NewHMap[int, hg.HSlice[int]]()

	for i := range iter.N(5) {
		md.Set(i, md.GetOrDefault(i, hg.NewHSlice[int]()).Append(i))
	}
}
```
</td>
</tr>
</table>

CopyDir copies the contents of the current directory to the destination directory.

<table>
<tr>
<th><code>stdlib</code></th>
<th><code>hg</code></th>
</tr>
<tr>
<td>

```go
func copyDir(src, dest string) error {
	return filepath.Walk(src, func(path string,
		info fs.FileInfo, err error,
	) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath, info.Mode())
	})
}

func copyFile(src, dest string, mode fs.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	return err
}

func main() {
	src := "path/to/source/directory"
	dest := "path/to/destination/directory"

	err := copyDir(src, dest)
	if err != nil {
		fmt.Println("Error copying directory:", err)
	} else {
		fmt.Println("Directory copied successfully")
	}
}
```
</td>
<td>

```go
func main() {
	d := hg.NewHDir(".").CopyDir("copy")

	if d.Error() != nil {
		fmt.Println(d.Error())
	}
}
```
</td>
</tr>
</table>

RandomSample returns a new slice containing a random sample of elements from the original slice.

<table>
<tr>
<th><code>stdlib</code></th>
<th><code>hg</code></th>
</tr>
<tr>
<td>

```go
func RandomSample(slice []int, amount int) []int {
	if amount > len(slice) {
		amount = len(slice)
	}

	samples := make([]int, amount)

	for i := 0; i < amount; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(slice))))
		samples[i] = slice[index.Int64()]
		slice = append(slice[:index.Int64()], slice[index.Int64()+1:]...)
	}

	return samples
}

func main() {
	slice := []int{1, 2, 3, 4, 5, 6}
	samples := RandomSample(slice, 3)
	fmt.Println(samples)
}
```
</td>
<td>

```go
func main() {
	slice := hg.HSliceOf(1, 2, 3, 4, 5, 6)
	samples := slice.RandomSample(3)
	fmt.Println(samples)
}
```
</td>
</tr>
</table>

Insert inserts values at the specified index in the slice and returns the resulting slice.

<table>
<tr>
<th><code>stdlib</code></th>
<th><code>hg</code></th>
</tr>
<tr>
<td>

```go
func Insert(slice []int, index int, values ...int) []int {
	total := len(slice) + len(values)
	if total <= cap(slice) {
		slice = slice[:total]
		copy(slice[index+len(values):], slice[index:])
		copy(slice[index:], values)

		return slice
	}

	newSlice := make([]int, total)
	copy(newSlice, slice[:index])
	copy(newSlice[index:], values)
	copy(newSlice[index+len(values):], slice[index:])

	return newSlice
}

func main() {
	slice := []int{1, 2, 3, 4, 5}
	slice = Insert(slice, 2, 6, 7, 8)
	fmt.Println(slice) // Output: [1 2 6 7 8 3 4 5]
}
```
</td>
<td>

```go
func main() {
	slice := hg.HSliceOf(1, 2, 3, 4, 5)
	slice = slice.Insert(2, 6, 7, 8)
	fmt.Println(slice) // Output: [1 2 6 7 8 3 4 5]
}
```
</td>
</tr>
</table>

Permutations returns all possible permutations of the elements in the slice.

<table>
<tr>
<th><code>stdlib</code></th>
<th><code>hg</code></th>
</tr>
<tr>
<td>

```go
func Permutations(slice []int) [][]int {
	if len(slice) <= 1 {
		return [][]int{slice}
	}

	perms := make([][]int, 0)

	for i, elem := range slice {
		rest := make([]int, len(slice)-1)

		copy(rest[:i], slice[:i])
		copy(rest[i:], slice[i+1:])

		subPerms := Permutations(rest)

		for j := range subPerms {
			subPerms[j] = append([]int{elem}, subPerms[j]...)
		}

		perms = append(perms, subPerms...)
	}

	return perms
}

func main() {
	slice := []int{1, 2, 3}
	fmt.Println(Permutations(slice))
	// Output: [[1 2 3] [1 3 2] [2 1 3] [2 3 1] [3 2 1] [3 1 2]]
}
```
</td>
<td>

```go
func main() {
	slice := hg.HSliceOf(1, 2, 3)
	fmt.Println(slice.Permutations())
	// Output: [[1 2 3] [1 3 2] [2 1 3] [2 3 1] [3 2 1] [3 1 2]]
}
```
</td>
</tr>
</table>
<br>
<p align="center">
  <img src="https://user-images.githubusercontent.com/65846651/233453773-33f38b64-0adc-41b4-8e13-a49c89bf9db6.png">
</p>

# 🤖👋 Human Surf: makes HTTP fun and easy!
Surf is a fun, user-friendly, and lightweight Go library that allows you to interact with HTTP services as if you were chatting with them face-to-face! 😄
Imagine if you could make HTTP requests by simply asking a server politely, and receiving responses as if you were having a delightful conversation with a friend. That's the essence of surf!

## 🌟 Features
1. 💬 **Simple and expressive:** surf API is designed to make your code look like a conversation, making it easier to read and understand.
2. 🚀 **Asynchronous magic:** With surf built-in async capabilities, you can send multiple requests in parallel and handle them effortlessly.
3. 💾 **Caching and streaming:** Efficiently cache response bodies and stream data on the fly, like a superhero saving the world from slow internet connections.
4. 📉 **Limit and deflate:** Limit the amount of data you receive, and decompress it on the fly, giving you more control over your HTTP interactions.
5. 🎩 **Flexible:** Customize headers, query parameters, timeouts, and more, for a truly tailor-made experience.

## 💻 Example
Here's a fun and friendly example of how surf makes HTTP requests look like a conversation:
```Go
package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	resp, err := surf.NewClient().Get("https://api.example.com/jokes/random").Do() // A simple GET request
	if err != nil { log.Fatal(err) }

	joke := struct {
		ID     int    `json:"id"`
		Setup  string `json:"setup"`
		Punch  string `json:"punch"`
	}{}

	resp.Body.JSON(&joke)

	fmt.Println("Joke of the day:")
	fmt.Printf("%s\n%s\n", joke.Setup, joke.Punch)
}
```

## 🚀 Getting Started
To start making friends with HTTP services, follow these simple steps:
1. Install the surf package using **go get:**
```bash
go get -u github.com/x0xO/hg
```
2. Import the package into your project:
```Go
import "github.com/x0xO/hg/surf"
```
3. Start making requests and have fun! 😄

Give surf a try, and watch your HTTP conversations come to life!
