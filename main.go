package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/mileusna/crontab"
	"net/http"
)

var databaseURL = "mongo"
var sslBit = true

func main() {

	PrepareTemplate()

	DialConnection(databaseURL)

	defer session.Close()

	ctab := crontab.New()
	ctab.AddJob("10 23 * * *", RefreshImages)

	http.HandleFunc("/github/", Flair)

	box := packr.NewBox("./secrets")

	certBundle, _ := box.MustBytes("crt-bundle.pem")

	sslPrivate, _ := box.MustBytes("ssl-private.key")

	certsPair, err := tls.X509KeyPair(certBundle, sslPrivate)

	if err == nil && sslBit {
		fmt.Println("Starting server with SSL on port 8443")
		config := &tls.Config{Certificates: []tls.Certificate{certsPair},
			// Turn off warning about self signed cert
			InsecureSkipVerify: true,
		}
		server := http.Server{
			TLSConfig: config,
			Addr:      ":8443",
		}
		server.ListenAndServeTLS("", "")

	} else {

		fmt.Println("Starting server without SSL on port 8443")
		http.ListenAndServe(":8443", nil)
	}

}
