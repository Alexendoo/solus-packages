package main

import (
	"os"

	"github.com/Alexendoo/solus-packages/packages"
)

func main() {
	// log.Println(packages.Download(packages.CachePath()))

	reader, _ := os.Open(`eopkg-index.xml`)
	packages.Decode(reader)

	// search.Search()
}
