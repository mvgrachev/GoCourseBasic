package main

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"log"
)

func main() {
	var dat map[string]interface{}
	resp, err := http.Get("https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key=https://disk.yandex.ru/i/xV3LGCEwax8RfA")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	
	if err = json.Unmarshal(b, &dat); err != nil {
		panic(err)
	}

	href := dat["href"].(string)	
	resp, err = http.Get(href)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err != nil {
		return
	}
	
	b, err = ioutil.ReadAll(resp.Body)

	err = ioutil.WriteFile("/Users/mvgrachev/go/src/GoCourseBasic/homework-1/testfile", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
