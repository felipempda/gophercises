package link

import (
	"golang.org/x/net/html"
	"io"
	"os"
	"strings"
)

type HrefLink struct {
	Href string
	Text string
}

func ParseReader(reader io.Reader) ([]HrefLink, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	return extractLinks(node), nil
}

// you can't overload functions with the same name and different argument types
// maybe it's for the best?
func ParseFile(filename string) ([]HrefLink, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return ParseReader(f)
}

func extractLinks(node *html.Node) []HrefLink {
	var h []HrefLink
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {

			for _, a := range n.Attr {
				if a.Key == "href" {
					text := findText(n)
					newLink := HrefLink{a.Val, text}
					h = append(h, newLink)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	return h
}

// ok so you need to run the siblings of N to find a suitable TextNode
func findText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += findText(c)
	}
	// remove space between words, better than trim (but more expensive)
	return strings.Join(strings.Fields(ret), " ")
}
