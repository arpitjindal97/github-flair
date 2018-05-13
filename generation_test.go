package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

// TestCreateFlair tests if flair is there is any error
// while creating the flair
func TestCreateFlair(t *testing.T) {

	err := PrepareTemplate()
	if err != nil {
		t.Error(err)
	}

	RequestFlair("http://example.com/github/arpitjindal97.png", t)
	RequestFlair("http://example.com/github/arpitjindal97.png?theme=dark", t)

	DeleteFolder()
	RequestFlair("http://example.com/github/arpitjindal97.png", t)

	RefreshImages()

}

// RequestFlair request the flair from handler
func RequestFlair(url string, t *testing.T) {

	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	Flair(w, req)

	resp := w.Result()
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
