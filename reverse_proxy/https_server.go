package main

import(
	"net/url"
	"net/http"
	"net/http/httputil"
)

func main() {


	remote1, err := url.Parse("http://localhost:8081")
	if err != nil {
		panic(err)
	}

	proxy1 := httputil.NewSingleHostReverseProxy(remote1)
	http.HandleFunc("/", handler(proxy1))

	http.ListenAndServeTLS(":443", "certificate.pem",
		"private.key", nil)
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Println(r.URL)
		//r.URL.Path = "/"
		p.ServeHTTP(w, r)
	}
}