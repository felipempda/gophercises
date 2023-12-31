package urlshort_02

import (
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		for path, url := range pathsToUrls {
			if r.URL.Path == path {
				log.Printf("Redirected request to %q", url)
				http.Redirect(w, r, url, http.StatusSeeOther)
				return
			}
		}

		log.Printf("Path not found: %q, using fallback\n", r.URL.Path)
		fallback.ServeHTTP(w, r)
	}

	for path, url := range pathsToUrls {
		http.HandleFunc(path, f)
		log.Printf("Registered path %q, redirect to %q\n", path, url)
	}

	return f
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

type redirectYaml struct {
	Path string // needs to start with uppercase letter because of visibility rules! (you can use tags as well `json:"path"`)
	Url  string
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml := []redirectYaml{}

	err := yaml.Unmarshal(yml, &parsedYaml)
	if err != nil {
		log.Printf("Error parse YAML..")
	} else {
		log.Printf("YAML parse successful!")
	}
	log.Println("YAML content: ", parsedYaml)

	f := func(w http.ResponseWriter, r *http.Request) {

		for _, item := range parsedYaml {
			if r.URL.Path == item.Path {
				log.Printf("Redirected request to %q", item.Url)
				http.Redirect(w, r, item.Url, http.StatusSeeOther)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}

	for _, item := range parsedYaml {
		log.Printf("Registered path %q, redirect to URL %q", item.Path, item.Url)
		http.HandleFunc(item.Path, f)
	}

	return f, err
}
