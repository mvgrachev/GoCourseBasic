package main

import (
	"io/ioutil"
	"net/http"
	"bytes"
	"sync"
	"log"
	"os"
	"os/signal"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/go-chi/chi"
	"time"
)


const (
	COOKIE_KEY = "ssid"
)

func main() {
	stopchan := make(chan os.Signal)

	router := chi.NewRouter()
	router.Route("/sites", func(r chi.Router) {
		r.Post("/", parseHandle)
	})		
	router.Route("/cookie", func(r chi.Router) {
		r.Get("/", setCookie)
		r.Get("/{cookieName}", getCookie)
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	signal.Notify(stopchan, os.Interrupt, os.Kill)
	<-stopchan
	log.Print("gracefull shutdown")
}

func parseHandle(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		wr.WriteHeader(http.StatusMethodNotAllowed)
		wr.Header().Set("Content-Type", "application/json")
		wr.Write([]byte(`{"Error": "Method not allowed"}`))
		return
	}

	type queryData struct {
		Search string `json:"search"`
		Sites  []string `json:"sites"`
	}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(req.Body)
	jsonBlob := buffer.Bytes()

	if json.Valid(jsonBlob) == false {
		wr.WriteHeader(http.StatusBadRequest)
		wr.Header().Set("Content-Type", "application/json")
		wr.Write([]byte(`{"Error": "Invalid json"}`))
		return
	}

	var query queryData
	err := json.Unmarshal(jsonBlob, &query)
	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		wr.Header().Set("Content-Type", "application/json")
		wr.Write([]byte(`{"Error": "Invalid encode json"}`))
		return
	}

	var res []string
	
	res, err = parse(query.Sites, query.Search)

	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		wr.Header().Set("Content-Type", "application/json")
		wr.Write([]byte(`{"Error": "Invalid request site"}`))
		return
	}

	type responseData struct {
		Sites []string `json:"sites"`
	}
	
	var response responseData
	response.Sites = res
	responseBlob, err := json.Marshal(response)
	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		wr.Header().Set("Content-Type", "application/json")
		wr.Write([]byte(`{"Error": "Invalid decode json"}`))
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.Write(responseBlob)
}

func setCookie(wr http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie(COOKIE_KEY)
	if cookie == nil {
		cookie = &http.Cookie{
			Name: COOKIE_KEY,
			Path: "/",
		}
	}

	cookie.Value = uuid.Must(uuid.NewV4()).String()
	expire := time.Now().AddDate(0, 0, 1)
	cookie.Expires = expire
	http.SetCookie(wr, cookie)

	responseBlob, err := json.Marshal(cookie)
	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		wr.Header().Set("Content-Type", "application/json")
		wr.Write([]byte(`{"Error": "Invalid decode json"}`))
		return
	}
	wr.Write(responseBlob)
}

func getCookie(wr http.ResponseWriter, req *http.Request) {
	cookieName := chi.URLParam(req, "cookieName")
	
	cookie, _ := req.Cookie(cookieName)

	if cookie == nil {
		wr.Header().Set("Content-Type", "application/json")
		wr.Write([]byte(`{"Cookie": "Undef"}`))
		return
	}

	responseBlob, err := json.Marshal(cookie)
	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		wr.Header().Set("Content-Type", "application/json")
		wr.Write([]byte(`{"Error": "Invalid decode json"}`))
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	wr.Write(responseBlob)
}

func parse(urls []string, query string) (a []string, err error) {
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
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
	wg.Wait()
	return
}
