package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/felipempda/gophercises/13_quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	storiesCached := storyCache{
		numStories: numStories,
		duration:   6 * time.Second,
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			temp := storyCache{
				numStories: numStories,
				duration:   6 * time.Second,
			}
			temp.stories()
			storiesCached.mutex.Lock()
			storiesCached.cache = temp.cache
			storiesCached.duration = temp.duration
			storiesCached.expiration = time.Now().Add(storiesCached.duration)
			storiesCached.mutex.Unlock()
			<-ticker.C
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := storiesCached.stories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type storyCache struct {
	cache      []item
	expiration time.Time
	duration   time.Duration
	mutex      sync.Mutex
	numStories int
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if time.Now().Sub(sc.expiration) < 0 {
		fmt.Println("cache hit")
		return sc.cache, nil
	}
	fmt.Println("cache miss")
	items, err := getTopStories(sc.numStories)
	if err != nil {
		return nil, err
	}
	sc.expiration = time.Now().Add(sc.duration)
	sc.cache = items
	return sc.cache, nil
}

var (
	cache           []item
	cacheExpiration time.Time
	cacheMutex      sync.Mutex
)

// could Not simulate Race
func getCachedTopStories(numStories int) ([]item, error) {

	// expired Cache
	var items []item
	var err error
	if time.Now().Sub(cacheExpiration) > 0 || len(cache) == 0 {
		fmt.Println("cache miss")
		items, err = getTopStories(numStories)
		cache = items
		cacheExpiration = time.Now().Add(10 * time.Second)
	} else {
		fmt.Println("cache hit")
		items = cache
		err = nil
	}
	return items, err
}

// could Not simulate Race
func getCachedTopStoriesProf(numStories int) ([]item, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if time.Now().Sub(cacheExpiration) < 0 {
		return cache, nil
	}

	items, err := getTopStories(numStories)
	if err != nil {
		return nil, err
	}
	cache = items
	cacheExpiration = time.Now().Add(10 * time.Second)
	return cache, nil
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}
	var resultItems []item
	at := 0
	for len(resultItems) < numStories {
		needed := (numStories - len(resultItems)) * 5 / 4 // add 25%
		temp := getStories(ids[at:at+needed], client)
		fmt.Printf("loop at=%d, needed=%d, numStories=%d found=%d\n", at, needed, numStories, len(temp))
		resultItems = append(resultItems, temp...)
		at += needed
	}
	return resultItems[:numStories], nil
}

func getStories(ids []int, client hn.Client) []item {
	type result struct {
		idx  int
		item item
		err  error
	}
	resultCh := make(chan result)
	for i := 0; i < len(ids); i++ {
		go func(idx, id int) {
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{idx: idx, err: err}
			}
			resultCh <- result{idx: idx, item: parseHNItem(hnItem)}
		}(i, ids[i])
	}
	resultItems := make([]result, 0, 0)
	for i := 0; i < len(ids); i++ {
		resultItems = append(resultItems, <-resultCh)
	}
	sort.Slice(resultItems, func(i, j int) bool {
		return resultItems[i].idx < resultItems[j].idx
	})
	var stories []item
	for _, res := range resultItems {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}
	return stories
}
func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
