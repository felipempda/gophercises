package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting webserver...")
	cyoa := readJson()
	http.ListenAndServe(":8080", cyoaHandler(cyoa))
}

func readJson() Cyoa {
	data, err := os.ReadFile("./gopher.json")
	if err != nil {
		panic(err)
	}
	cyoa := parseJson(data)
	return cyoa
}

type Arch struct {
	Title   string
	Story   []string
	Options []Option
}

type Option struct {
	Text string
	Arc  string
}

type Cyoa map[string]Arch

func parseJson(bytes []byte) Cyoa {
	cyoa := Cyoa{}
	if err := json.Unmarshal(bytes, &cyoa); err != nil {
		panic(err)
	}
	log.Println("Json parsed!")
	// log.Println(cyoa["debate"])
	// log.Println(cyoa["debate"].Title)
	// log.Println(cyoa["debate"].Story[0])
	// log.Println(cyoa["debate"].Options[0].Arc)
	return cyoa
}

func cyoaHandler(cyoa Cyoa) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:] //remove initial /
		if arch, ok := cyoa[path]; ok {
			fmt.Fprintf(w, "<p>You are at Arch %q (%q)</p><br>", path, arch.Title)
			for _, title := range arch.Story {
				fmt.Fprintln(w, title)
			}
		} else {
			fmt.Fprintf(w, "Arch %q not found", path)
		}
	})
}
