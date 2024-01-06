package sitemap

import (
	"fmt"
	link "github.com/felipempda/gophercises/04_link/link"
	"net/http"
	"strings"
)

func Sitemap(website string) {
	resp, err := http.Get(website)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	links, err := link.ParseReader(resp.Body)
	for _, link := range links {
		normalizedLink, sameWebsite := normalizeLink(link.Href, website)
		fmt.Printf("Href: %v, Link: %v, Normalized: %v, SameWebiste: %v\n", link.Href, link.Text, normalizedLink, sameWebsite)
	}
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
