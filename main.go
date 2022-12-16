package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Domains struct {
	Domain string `json:"domain"`
	Proxy  string `json:"proxy"`
	SSL    bool   `json:"ssl"`
}

func main() {

	// 	//Read JSON file
	jsonFile, err := os.ReadFile("dauqu.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Opened dauqu.json")

	var dauqu []Domains

	// 	//Unmarshal JSON file
	err = json.Unmarshal(jsonFile, &dauqu)
	if err != nil {
		log.Fatal(err)
	}

	//Loop through domains
	for _, domain := range dauqu {

		vhost, err := url.Parse(domain.Proxy)
		if err != nil {
			log.Fatal(err)
		}

		//Create proxy
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: vhost.Scheme,
			Host:   vhost.Host,
		})

		//Set Header
		// proxy.Director = func(req *http.Request) {
		// 	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		// 	req.Header.Set("Host", vhost.Host)
		// 	req.URL.Scheme = vhost.Scheme
		// 	req.URL.Host = vhost.Host
		// }

		//Header response
		// proxy.ModifyResponse = func(resp *http.Response) error {
		// 	resp.Header.Set("Server", "Setkaro")
		// 	resp.Header.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		// 	resp.Header.Set("Alt-Svc", "h2=\":443\"; ma=2592000")
		// 	resp.Header.Set("X-Forwarded-Proto", "https")
		// 	return nil
		// }

		http.HandleFunc(domain.Domain+"/", func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
	}

	//Listen and serve
	http.ListenAndServe(":80", nil)
}
