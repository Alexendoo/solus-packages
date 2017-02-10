package blevesearch

import (
	"fmt"

	"os"

	"log"

	"github.com/Alexendoo/solus-packages/packages"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
)

func Index(packages *packages.PISI, index bleve.Index) {
	for _, obsolete := range packages.Distribution.Obsoletes {
		err := index.Delete(obsolete)
		if err != nil {
			log.Println(err)
		}
	}

	for _, pkg := range packages.Packages {
		err := index.Index(pkg.Name, pkg)
		if err != nil {
			log.Println(err)
		}
	}
}

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

func getMapping() mapping.IndexMapping {
	mapping := bleve.NewIndexMapping()
	mapping.DefaultAnalyzer = en.AnalyzerName

	return mapping
}

func getIndex(path string) (bleve.Index, error) {
	index, err := bleve.Open(path)
	if err == nil {
		return index, nil
	}
	if err != bleve.ErrorIndexPathDoesNotExist {
		return nil, err
	}

	mapping := getMapping()

	index, err = bleve.New(path, mapping)
	if err != nil {
		return nil, err
	}

	return index, nil
}
