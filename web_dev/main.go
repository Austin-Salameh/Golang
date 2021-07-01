package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Sitemap struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles            []string `xml:"url>news>title"`
	Keywords          []string `xml:"url>news>keywords"`
	Locations         []string `xml:"url>loc"`
	Publication_Dates []string `xml:"url>news>publication_date"`
}

type NewsMap struct {
	Keyword          string
	Location         string
	Publication_Date string
}

type NewsAggPage struct {
	Title string
	News  map[string]NewsMap
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var s Sitemap
	var n News
	resp, _ := http.Get("https://www.cbssports.com/index-sitemap.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
	news_map := make(map[string]NewsMap)
	// fmt.Println(s.Locations)

	for _, Location := range s.Locations {
		resp, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(bytes, &n)

		for i, _ := range n.Keywords {
			news_map[n.Titles[i]] = NewsMap{n.Keywords[i], n.Locations[i], n.Publication_Dates[i]}
		}
	}

	p := NewsAggPage{Title: "Sports News Aggregator", News: news_map}
	t, err := template.ParseFiles("displayNews.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, p)
}

func main() {

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
