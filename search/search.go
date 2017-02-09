package blevesearch

import (
	"fmt"

	"os"

	"github.com/blevesearch/bleve"
)

func Search() {
	// open a new index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("example.bleve", mapping)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := struct {
		Name string
	}{
		Name: "foo bar",
	}

	data2 := struct {
		Name string
	}{
		Name: "foo baz",
	}

	// index some data
	index.Index("text", data)
	index.Index("foo", data2)

	// search for some text
	queryOne := bleve.NewMatchQuery("foo")
	queryTwo := bleve.NewMatchQuery("bar")
	queryThree := bleve.NewMatchQuery("baz")

	query := bleve.NewDisjunctionQuery(queryOne, queryTwo, queryThree)

	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)

	index.Close()
	os.RemoveAll("example.bleve")
}
