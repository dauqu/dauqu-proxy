package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

type Domains struct {
	Domain string `json:"domain"`
	Proxy  string `json:"proxy"`
	SSL    bool   `json:"ssl"`
}

func main() {

	//Read JSON file
	jsonFile, err := os.ReadFile("dauqu.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Opened dauqu.json")

	var dauqu []Domains

	//Unmarshal JSON file
	err = json.Unmarshal(jsonFile, &dauqu)
	if err != nil {
		log.Fatal(err)
	}

	//Lpp through domains
	for _, domain := range dauqu {

		//Proxy URL
		vhost, err := url.Parse(domain.Proxy)
		if err != nil {
			panic(err)
		}

		//Create proxy
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   domain.Proxy,
		})

		//Headers
		proxy.Director = func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", vhost.Host)
			req.URL.Scheme = vhost.Scheme
			req.URL.Host = vhost.Host
		}

		http.HandleFunc(domain.Domain, handler(proxy))

		//Autocert
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain.Domain),
			Cache:      autocert.DirCache("certs"),
		}

		//Create TLS config
		tlsConfig := &tls.Config{
			GetCertificate: m.GetCertificate,
		}

		//Create server
		server := &http.Server{
			Addr:      ":443",
			TLSConfig: tlsConfig,
		}

		//Run server
		go server.ListenAndServeTLS("", "")
	}

	//Listen on port 80
	go http.ListenAndServe(":80", nil)
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter,
	*http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}

// func main() {

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello Secure World")
// 	})

// 	//Create autocert manager
// 	m := autocert.Manager{
// 		Prompt:     autocert.AcceptTOS,
// 		HostPolicy: autocert.HostWhitelist("a.setkaro.com"),
// 		Cache:      autocert.DirCache("certs"),
// 	}

// 	//Create TLS config
// 	tlsConfig := &tls.Config{
// 		GetCertificate: m.GetCertificate,
// 	}

// 	//Create server
// 	s := &http.Server{
// 		Addr:      ":443",
// 		TLSConfig: tlsConfig,
// 		Handler:   mux,
// 	}

// 	//Listen on port 80
// 	go http.ListenAndServe(":80", nil)

// 	//Listen on port 443
// 	s.ListenAndServeTLS("", "")
// }
