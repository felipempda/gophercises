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
		panic(err)
	}

	node, err := html.Parse(f)
	if err != nil {
		panic(err)
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
					newLink := HrefLink{a.Val, strings.TrimSpace(n.FirstChild.Data)} // how to find this bloody text ?
					h.HrefLinks = append(h.HrefLinks, newLink)
					break
				}
			}
		}
		// } else if n.Type == html.TextNode {
		//    log.Printf("Text = %v", n.Data)
		// }
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
}
