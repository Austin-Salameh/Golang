package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

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

func newsRoutine(c chan News, Location string) {
	defer wg.Done()
	var n News
	resp, _ := http.Get(Location)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &n)
	resp.Body.Close()
	c <- n
}

func queueWatcher(c chan News) map[string]NewsMap {
	news_map := make(map[string]NewsMap)

	for elem := range c {
		for i, _ := range elem.Keywords {
			news_map[elem.Titles[i]] = NewsMap{elem.Keywords[i], elem.Locations[i], elem.Publication_Dates[i]}
		}
	}

	return news_map
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var s Sitemap

	resp, _ := http.Get("https://www.cbssports.com/index-sitemap.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
	var news_map map[string]NewsMap
	resp.Body.Close()

	queue := make(chan News, 30)
	processingFinished := make(chan bool)

	go func() {
		news_map = queueWatcher(queue)
		processingFinished <- true
	}()

	for _, Location := range s.Locations {
		wg.Add(1)
		go newsRoutine(queue, Location)
	}

	wg.Wait()
	close(queue)
	<-processingFinished

	fmt.Println(news_map)

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
