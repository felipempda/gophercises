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
	links := getLinks(website, maxLevel)
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

type siteMap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []url    `xml:"url"`
}

type url struct {
	Link string `xml:"loc"`
}

func buildXML(links []string) []byte {

	u := siteMap{}
	u.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9" // this feels illegal
	for _, link := range links {
		u.Urls = append(u.Urls, url{link})
	}
	var result []byte
	result, err := xml.MarshalIndent(u, "", " ")
	if err != nil {
		panic(err)
	}
	var xmlHeader = []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n") // this as well lol
	return append(xmlHeader, result...)

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
