package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/autotls"
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

	//Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// 	//Read JSON file
	jsonFile, err := os.ReadFile("/var/dauqu/dauqu.json")
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

	//Add default domain to dauqu
	dauqu = append(dauqu, Domains{
		Domain: hostname,
		Proxy:  "http://localhost:9000",
		SSL:    true,
	})

	mux := http.NewServeMux()

	var domains []string

	//Loop through domains
	for _, domain := range dauqu {

		//Add domain to domains
		domains = append(domains, domain.Domain)

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
		proxy.Director = func(req *http.Request) {
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
			req.Header.Set("Host", vhost.Host)
			req.URL.Scheme = vhost.Scheme
			req.URL.Host = vhost.Host
		}

		//Header response
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Set("Server", "Setkaro")
			resp.Header.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
			resp.Header.Set("Alt-Svc", "h2=\":443\"; ma=2592000")
			resp.Header.Set("X-Forwarded-Proto", "https")
			return nil
		}

		mux.HandleFunc(domain.Domain+"/", func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
	}

	//Add hostname to domains
	domains = append(domains, hostname)

	//Listen and serve
	autotls.Run(mux, domains...)
}
