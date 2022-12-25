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
		//If file not found create new file
		if os.IsNotExist(err) {
			_, err = os.Create("/var/dauqu/dauqu.json")
			if err != nil {
				log.Fatal(err)
			}

			//Write to file
			err = os.WriteFile("/var/dauqu/dauqu.json", []byte("[]"), 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	fmt.Println("Successfully Opened dauqu.json")

	var dauqu []Domains

	// 	//Unmarshal JSON file
	err = json.Unmarshal(jsonFile, &dauqu)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	var domains []string

	//Loop through domains
	if len(dauqu) > 0 {
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
				resp.Header.Set("X-Forwarded-Proto", "https")
				return nil
			}

			mux.HandleFunc(domain.Domain+"/", func(w http.ResponseWriter, r *http.Request) {
				//if proxy not found show 404
				if vhost.Host == "" {
					//404.html
					http.ServeFile(w, r, "/var/dauqu/404.html")
					return
				} else {
					proxy.ServeHTTP(w, r)
				}
			})
		}
	}

	vhost, err := url.Parse("http://localhost:9000")
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

	//Allow to use cookie
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	mux.HandleFunc(hostname+"/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	//Add default domain
	domains = append(domains, hostname)

	//Listen and serve
	autotls.Run(mux, domains...)
}
