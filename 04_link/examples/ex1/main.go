package main

import (
	link "github.com/felipempda/gophercises/04_link/link"
	"log"
	"strings"
)

func main() {

	html := `
	<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>
`
	// testParseReader
	reader := strings.NewReader(html)
	links, err := link.ParseReader(reader)
	if err != nil {
		panic(err)
	}
	for _, link := range links {
		log.Printf("ref=%v , text=%v\n", link.Href, link.Text)
	}

	// testParseFile
	links, err = link.ParseFile("examples/ex1/ex1.html")
	if err != nil {
		panic(err)
	}
	for _, link := range links {
		log.Printf("ref=%v , text=%v\n", link.Href, link.Text)
	}
}
