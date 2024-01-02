package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting webserver...")
	cyoa := readJson()
	http.ListenAndServe(":8080", cyoa)
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

func (cyoa Cyoa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:] //remove initial /
	if arch, ok := cyoa[path]; ok {
		runArchTemplate(arch, w)
		return
	}
	if path == "" {
		http.Redirect(w, r, "/intro", http.StatusSeeOther)
		return
	}
	fmt.Fprintf(w, "Arch %q not found", path)
}

func runArchTemplate(arch Arch, w http.ResponseWriter) {
	t, err := template.New("webpage").ParseFiles("./arch.html")
	err = t.ExecuteTemplate(w, "T", arch)
	if err != nil {
		fmt.Fprintf(w, "Oops, error rendering template. Please be a better programmer!")
		log.Println(err)
	}
}
