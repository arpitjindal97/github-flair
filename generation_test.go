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

	// BothTheme("arpitjindal97", t)
	// BothTheme("narkoz", t)
	req := httptest.NewRequest("GET", "http://example.com/github/arpitjindal97.png", nil)
	w := httptest.NewRecorder()
	Flair(w, req)

	resp := w.Result()
	header := resp.Header.Get("Content-Type")

	if header != "image/png" {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Error(body)
	}
	RefreshImages()

}

// BothTheme tests both themes of given username
// saves me few lines of code
func BothTheme(username string, t *testing.T) {
	_, err := CreateFlair(username, "clean")
	if err != nil {
		t.Error(err)
	}

	_, err = CreateFlair(username, "dark")
	if err != nil {
		t.Error(err)
	}
}
