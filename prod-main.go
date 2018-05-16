// +build prod

package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/mileusna/crontab"
	"io/ioutil"
	"net/http"
)

func main() {

	PrepareTemplate()

	DialConnection("mongo")

	defer session.Close()

	ctab := crontab.New()
	ctab.AddJob("10 23 * * *", RefreshImages)

	http.HandleFunc("/github/", Flair)

	box := packr.NewBox("./secrets")

	certBundle, _ := box.MustBytes("crt-bundle.pem")

	sslPrivate, _ := box.MustBytes("ssl-private.key")

	certsPair, err := tls.X509KeyPair(certBundle, sslPrivate)
	if err != nil {
		fmt.Println("error in creating x509 pair", err)

	}

	config := &tls.Config{Certificates: []tls.Certificate{certsPair},
		// Turn off warning about self signed cert
		InsecureSkipVerify: true,
	}

	server := http.Server{
		TLSConfig: config,
		Addr:      ":443",
	}

	server.ListenAndServeTLS("", "")

}
