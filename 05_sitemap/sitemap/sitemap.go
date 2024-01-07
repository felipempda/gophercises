package sitemap

import (
	"encoding/xml"
	"fmt"
	link "github.com/felipempda/gophercises/04_link/link"
	"net/http"
	"strings"
)

// negative maxLevel means you don't care about them
func SitemapBuilder(website string, maxLevel int) []byte {
	links := bfs(website, maxLevel)
	// links :=  getLinks(website, maxLevel) // initial implementation with recursive functions
	xml := buildXML(links)
	return xml
}

func getLinks(website string, maxLevel int) []string {
	visitedPages := make(map[string]bool)

	var navigateLinks func(string, int)
	navigateLinks = func(page string, level int) {
		if visited, _ := visitedPages[page]; visited {
			fmt.Printf("Page already visited: %v!\n", page)
			return
		}
		if level > maxLevel && maxLevel >= 0 {
			return
		}
		fmt.Printf("At level %d\n", level)
		links := GetLinksFromSameWebsite(page, website)
		visitedPages[page] = true
		// mark new Links as not visited and visit them
		for _, link := range links {
			if visited, found := visitedPages[link]; !visited || !found {
				visitedPages[link] = false
				navigateLinks(link, level+1)
			}
		}
	}

	navigateLinks(website, 0)
	var links []string
	for page, _ := range visitedPages {
		links = append(links, page)
	}
	return links
}

func bfs(website string, maxLevel int) []string {
	visited := make(map[string]struct{})
	var q map[string]struct{}
	n := map[string]struct{}{
		website: struct{}{}, // {}{} here is because type=struct{} and initialization={}
		// otherwise you would have error => struct{} (type) is not an expression
	}
	for level := 0; level <= maxLevel || maxLevel < 0; level++ {
		fmt.Printf("At level %d\n", level)
		q, n = n, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for page, _ := range q {
			if _, ok := visited[page]; ok {
				continue
			}
			for _, link := range GetLinksFromSameWebsite(page, website) {
				if _, ok := visited[page]; !ok {
					n[link] = struct{}{}
				}
			}
			visited[page] = struct{}{}
		}
	}
	//   which one is better ?
	//
	links := make([]string, len(visited)) // len=size, cap=size
	var i = 0
	for link, _ := range visited {
		links[i] = link
		i++
	}

	// or
	//   links = make([]string, 0, len(visited)) // len=0, cap=size, so further appends only increment size while cap remains the same
	//   for link, _ := range visited {
	//   	links = append(links, link)
	//   }
	return links
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type siteMap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []url    `xml:"url"`
}

type url struct {
	Link string `xml:"loc"`
}

func buildXML(links []string) []byte {

	u := siteMap{
		Xmlns: xmlns,
	}
	for _, link := range links {
		u.Urls = append(u.Urls, url{link})
	}
	var result []byte
	result, err := xml.MarshalIndent(u, "", " ")
	panicThis(err)
	return append([]byte(xml.Header), result...)

}

func GetLinksFromSameWebsite(page, website string) []string {
	fmt.Printf("Visiting %v ...", page)

	resp, err := http.Get(page)
	panicThis(err)
	defer resp.Body.Close()

	links, err := link.ParseReader(resp.Body)
	panicThis(err)
	var linksList []string
	for _, link := range links {
		normalizedLink, sameWebsite := normalizeLink(link.Href, website)
		// fmt.Printf("Href: %v, Link: %v, Normalized: %v, SameWebiste: %v\n", link.Href, link.Text, normalizedLink, sameWebsite)
		if sameWebsite {
			linksList = append(linksList, normalizedLink)
		}
	}
	fmt.Printf("Found %d links\n", len(linksList))
	return linksList
}

func normalizeLink(link string, website string) (normalizedLink string, isFromWebsite bool) {
	if link[:1] == "/" {
		// relative path
		return website + link, true
	}
	if strings.HasPrefix(link, website) {
		return link, true
	}
	return link, false
}

func panicThis(error error) {
	if error != nil {
		panic(error)
	}
}
