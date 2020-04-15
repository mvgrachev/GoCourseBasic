package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"bytes"
	"flag"
	"sync"
)

func main() {
	var query string
	var urls string
	flag.StringVar(&query, "query", "query string", "Поисковый запрос")
	flag.StringVar(&urls, "urls", "https://ya.ru", "Список url через запятую без пробелов")
	
	flag.Parse()
	
	parse(urls, query)
}

func parse(urls string, query string) (a []string, err error) {
	var mutex = &sync.Mutex{}
	for _, url := range strings.Split(urls, ",") {
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}
			if ( bytes.Contains([]byte(b), []byte(query)) ) {
				mutex.Lock()
				a = append(a, url)
				mutex.Unlock()
			}
		}(url)
	}

	return
}
