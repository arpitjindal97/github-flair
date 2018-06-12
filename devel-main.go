// +build !prod

package main

import (
	"github.com/mileusna/crontab"
	"net/http"
)

var databaseURL = "mongo"

func main() {

	PrepareTemplate()

	DialConnection(databaseURL)

	defer session.Close()

	ctab := crontab.New()
	ctab.AddJob("10 23 * * *", RefreshImages)

	http.HandleFunc("/github/", Flair)

	// http.ListenAndServeTLS(":443", "secrets/crt-bundle.pem",
	//		"secrets/ssl-private.key", nil)

	http.ListenAndServe(":8080", nil)

}
