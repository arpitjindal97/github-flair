package main

import (
	"bytes"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Flair handles /github/ URL requests
// The URL should be  in this format
// github/username.png It extracts the username and finds it in
// database, If it exists then image is
// provided from folder else it os created
// and put in folder and returned
func Flair(w http.ResponseWriter, r *http.Request) {

	// extracts username from url
	username := strings.SplitN(r.URL.Path[1:], "/", 2)[1]
	username = username[:len(username)-4]

	theme := r.URL.Query().Get("theme")
	if theme == "" {
		theme = "clean"
	}

	log.Println(username)

	user := GetUser(username)

	if user.Username == "" {

		user.Username = username

		image, _ := CreateFlair(username, "clean")
		buffer := new(bytes.Buffer)
		png.Encode(buffer, image)
		user.Clean.Data = buffer.Bytes()

		image, _ = CreateFlair(username, "dark")
		buffer = new(bytes.Buffer)
		png.Encode(buffer, image)
		user.Dark.Data = buffer.Bytes()

		user.Timestamp = time.Now()

		InsertUser(user)
	}

	var byt []byte

	if byt = user.Dark.Data; theme != "dark" {
		byt = user.Clean.Data
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(byt)))
	w.Write(byt)

}
