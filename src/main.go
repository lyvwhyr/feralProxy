package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	xj "github.com/basgys/goxml2json"
)

func getURLResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func onCORSRequest(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		fmt.Println("url query param missing from url")
		return
	}

	resp, err := getURLResponse(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseBodyString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseBodyString)
}

func getRTNews(w http.ResponseWriter, r *http.Request) {
	rtNewsURL := "http://www.rt.com/rss"
	resp, err := getURLResponse(rtNewsURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	json, err := xj.Convert(resp.Body)
	if err != nil {
		panic("Failed to convert RT rss feed xml body")
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(json.Bytes())
}

func main() {
	http.HandleFunc("/news/rt", getRTNews)
	http.HandleFunc("/cors/", onCORSRequest) // set router
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
