package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

func main() {
	fmt.Println("Calling the API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "applicaiton/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseObject Response
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %v\n", responseObject)
}
