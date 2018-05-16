package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

// TestCreateFlair tests if flair is there is any error
// while creating the flair
func TestCreateFlair(t *testing.T) {

	databaseUrl = "localhost"
	go func() {
		main()
	}()

	time.Sleep(time.Second * 2)

	log.Println("Requesting arpitjindal97 clean flair")
	RequestFlair("http://localhost:8080/github/arpitjindal97.png", t)

	log.Println("Requesting narkoz dark flair for png image")
	RequestFlair("http://localhost:8080/github/narkoz.png?theme=dark", t)

	log.Println("Refreshing the images")
	RefreshImages()

	log.Println("Requesting arpitjindal97 clean flair")
	RequestFlair("http://localhost:8080/github/arpitjindal97.png", t)

}

// RequestFlair request the flair from handler
func RequestFlair(url string, t *testing.T) {

	resp, _ := http.Get(url)

	header := resp.Header.Get("Content-Type")

	if header != "image/png" {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Error(body)
	}
}
