package main

import (
	link "gophercises/04_link/link"
	"log"
	"os"
)

func main() {
	filename := os.Args[1]
	links, err := link.Parse(filename)
	if err != nil {
		panic(err)
	}
	log.Printf("len = %d\n", len(links.HrefLinks))
	for _, alink := range links.HrefLinks {
		log.Printf("ref=%v , text=%v\n", alink.Href, alink.Text)
	}
}
