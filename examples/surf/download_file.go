package main

import (
	"log"
	"net/url"
	"path"

	"github.com/x0xO/hg/surf"
)

func main() {
	dURL := "http://download.geonames.org/export/dump/alternateNames.zip"

	r, err := surf.NewClient().Get(dURL).Do()
	if err != nil {
		log.Fatal(err)
	}

	URL, err := url.ParseRequestURI(dURL)
	if err != nil {
		log.Fatal(err)
	}

	r.Body.Dump(path.Base(URL.Path))

	// or
	// r.Body.Dump("/home/user/some_file.zip")
}
