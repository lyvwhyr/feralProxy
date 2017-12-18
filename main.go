package rtrss

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func getURLResponse(url string, r *http.Request) (*http.Response, error) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	resp, err := client.Get(url)
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

	resp, err := getURLResponse(url, r)
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

type RTRSS struct {
	RSSItems []RSSItem `xml:"channel>item" json:"articles"`
}

type RSSItem struct {
	Title       string `xml:"title" json:"title"`
	GUID        string `xml:"guid" json:"guid"`
	Description string `xml:"description" json:"description"`
	PubDate     string `xml:"pubDate" json:"pubDate"`
}

func getRTNews(w http.ResponseWriter, r *http.Request) {
	rtNewsURL := "http://www.rt.com/rss"
	resp, err := getURLResponse(rtNewsURL, r)
	if err != nil {
		log.Println(err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var rtRss RTRSS
	err = xml.Unmarshal(bodyBytes, &rtRss)
	if err != nil {
		log.Println("Failed to convert RT rss feed xml body")
		log.Fatalln(err)
		return
	}
	fmt.Printf("%#v\n", rtRss)
	jsonData, err := json.Marshal(rtRss.RSSItems)
	if err != nil {
		log.Println("Failed to marshal RT rss model")
		log.Fatalln(err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsonData)
}

func init() {
	http.HandleFunc("/titan/news/rt", getRTNews)
	// http.HandleFunc("/titan/cors/", onCORSRequest) // set router
}

// <rss xmlns:media="http://search.yahoo.com/mrss/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:dc="http://purl.org/dc/elements/1.1/" version="2.0">
// 	<channel>
// 		<title>RT - Daily news</title>
// 		<link>https://www.rt.com</link>
// 		<description>RT : Today</description>
// 		<language>en</language>
// 		<copyright>RT</copyright>
// 		<atom:link href="https://www.rt.com/rss" rel="self" type="application/rss+xml"/>
// 		<image>
// 			<url>https://www.rt.com/static/img/logo-rss.png</url>
// 			<title>RT - Daily news</title>
// 			<link>https://www.rt.com</link>
// 			<width>125</width>
// 			<height>40</height>
// 		</image>
// 		<item>
// 			<title>
// 			Das boot ist kaputt! Germany has the world's best submarines... but none of them work
// 			</title>
// 			<link>
// 				<![CDATA[
// 				https://www.rt.com/news/413483-germany-submarines-out-of-service/?utm_source=rss&utm_medium=rss&utm_campaign=RSS
// 				]]>
// 			</link>
// 			<guid>
// 				https://www.rt.com/news/413483-germany-submarines-out-of-service/
// 			</guid>
// 			<description>
// 				<![CDATA[
// 				<img alt="Preview" align="left" style="margin-right: 10px;" src="https://cdni.rt.com/files/2017.12/thumbnail/5a36e2d5fc7e93a1128b4569.jpg" /> Germany is effectively without its entire submarine fleet, and won't have one vessel operational for months to come. Each one of the navy's vaunted U-boats is either on maintenance or in desperate need of repairs. <br/><a href="https://www.rt.com/news/413483-germany-submarines-out-of-service/?utm_source=rss&utm_medium=rss&utm_campaign=RSS">Read Full Article at RT.com</a>
// 				]]>
// 			</description>
// 			<enclosure url="https://cdni.rt.com/files/2017.12/thumbnail/5a36e2d5fc7e93a1128b4569.jpg" type="image/jpeg" length="123"/>
// 			<pubDate>Mon, 18 Dec 2017 04:34:49 +0000</pubDate>
// 			<dc:creator>RT</dc:creator>
// 		</item>
// 	</channel>
// </rss>
