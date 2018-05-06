// +build !devel

package main

import (
	"github.com/mileusna/crontab"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	PrepareTemplate()

	_, err := ioutil.ReadDir("flairs")
	if err != nil {
		os.Mkdir("flairs", 0755)
	}

	ctab := crontab.New()
	ctab.AddJob("10 23 * * *", RefreshImages)

	http.HandleFunc("/github/", Flair)

	http.ListenAndServeTLS(":443", "secrets/crt-bundle.pem",
		"secrets/ssl-private.key", nil)

	// http.ListenAndServe(":8080", nil)

}
