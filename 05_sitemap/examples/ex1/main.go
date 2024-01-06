package main

import (
	"flag"
	"fmt"
	"github.com/felipempda/gophercises/05_sitemap/sitemap"
)

func main() {
	website := flag.String("website", "https://wwww.google.com", "Website to build sitemap")
	maxLevel := flag.Int("maxLevel", 2, "MaxLevel of hierarchy to search links")
	flag.Parse()
	result := sitemap.SitemapBuilder(*website, *maxLevel)
	fmt.Println(string(result))
}
