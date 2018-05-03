package main

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	DownloadImages()

	http.HandleFunc("/github/", Flair)

	http.ListenAndServeTLS(":8080", "certificate.pem",
		"ssl-private.key", nil)

	// http.ListenAndServe(":8080", nil)

}

func Flair(w http.ResponseWriter, r *http.Request) {

	username := strings.SplitN(r.URL.Path[1:], "/", 2)[1]
	username = username[:len(username)-4]

	theme := r.URL.Query().Get("theme")
	if theme == "" {
		theme = "clean"
	}

	var myimage image.Image

	// if ExistsInDatabase(username) == true {
	//	myimage = GetFromFolder(username, theme)
	// } else {
	myimage = CreateFlair(username, theme)
	//	UpdateDatabase(username, myimage)
	//}

	w.Header().Set("Content-Type", "image/png")
	buffer := new(bytes.Buffer)
	png.Encode(buffer, myimage)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Write(buffer.Bytes())
}
