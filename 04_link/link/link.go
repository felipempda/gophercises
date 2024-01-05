package link

import (
	"golang.org/x/net/html"
	"os"
	"strings"
)

type HtmlParsed struct {
	HrefLinks []HrefLink
}

type HrefLink struct {
	Href string
	Text string
}

func Parse(filename string) (HtmlParsed, error) {
	h := HtmlParsed{}

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return h, err
	}

	node, err := html.Parse(f)
	if err != nil {
		return h, err
	}
	h.fillNode(node)

	return h, err
}

func (h *HtmlParsed) fillNode(node *html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {

			for _, a := range n.Attr {
				if a.Key == "href" {
					text := text(n)
					newLink := HrefLink{a.Val, text}
					h.HrefLinks = append(h.HrefLinks, newLink)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
}

// ok so you need to run the siblings of N to find a suitable TextNode
func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}
