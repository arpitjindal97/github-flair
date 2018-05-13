package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

// TestCreateFlair tests if flair is there is any error
// while creating the flair
func TestCreateFlair(t *testing.T) {

	go func() {
		main()
	}()

	time.Sleep(time.Second * 2)

	log.Println("Requesting arpitjindal97 clean flair")
	RequestFlair("http://localhost:8080/github/arpitjindal97.png", t)

	log.Println("Requesting arpitjindal97 dark flair")
	RequestFlair("http://localhost:8080/github/arpitjindal97.png?theme=dark", t)

	log.Println("Removing folder")
	DeleteFolder()

	log.Println("Requesting arpitjindal97 clean flair")
	RequestFlair("http://localhost:8080/github/arpitjindal97.png", t)

	log.Println("Refreshing the images")
	RefreshImages()

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

// DeleteFolder deletes the folder containing
// all the flair images
func DeleteFolder() {
	os.RemoveAll("data-db")
}
