package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getURLResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func onRequest(w http.ResponseWriter, r *http.Request) {
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

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseBodyBytes)
}

func main() {
	http.HandleFunc("/news/rt", getRTNews)
	http.HandleFunc("/cors/", onRequest)     // set router
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
