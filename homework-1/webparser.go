package main

import (
	"io/ioutil"
	"net/http"
	"GoCourseBasic/homework-1/webparse"
	"strings"
	"bytes"
)

func main() {
	var query string
	var urls string
	webparse.GetQuery(&query)
	webparse.GetUrls(&urls)
	
	res, _ := parse(urls, query)
	webparse.SayUrls(res)
}

func parse(urls string, query string) (a []string, err error) {
	for _, url := range strings.Split(urls, ",") {
		resp, err := http.Get(url)
		if err != nil {
			return a, err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if ( bytes.Contains([]byte(b), []byte(query)) ) {
			a = append(a, url)
		}
	}

	return a, err
}
