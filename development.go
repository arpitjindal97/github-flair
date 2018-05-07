// +build devel

package main

import (
	"github.com/mileusna/crontab"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	PrepareTemplate()

	_, err := ioutil.ReadDir("/data/flair-images")
	if err != nil {
		os.Mkdir("/data/flair-images", 0755)
	}

	ctab := crontab.New()
	ctab.AddJob("10 23 * * *", RefreshImages)

	http.HandleFunc("/github/", Flair)

	// http.ListenAndServeTLS(":443", "secrets/crt-bundle.pem",
	//		"secrets/ssl-private.key", nil)

	http.ListenAndServe(":8080", nil)

}
