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
	resp, err := getURLResponse(url)
	if err != nil {
		fmt.Println(err)
	}

	responseBodyString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseBodyString)
	w.WriteHeader(200)

}

func main() {
	http.HandleFunc("/", onRequest)          // set router
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
