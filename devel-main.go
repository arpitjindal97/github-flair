// +build !prod

package main

import (
	"github.com/mileusna/crontab"
	"net/http"
)

func main() {

	PrepareTemplate()

	CreateFolder()

	ctab := crontab.New()
	ctab.AddJob("10 23 * * *", RefreshImages)

	http.HandleFunc("/github/", Flair)

	// http.ListenAndServeTLS(":443", "secrets/crt-bundle.pem",
	//		"secrets/ssl-private.key", nil)

	http.ListenAndServe(":8080", nil)

}
