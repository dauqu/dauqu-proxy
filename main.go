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

	"github.com/gin-gonic/autotls"
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

			proxy.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}

			//Set Header
			proxy.Director = func(req *http.Request) {
				req.URL.Scheme = vhost.Scheme
				req.URL.Host = vhost.Host
				req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
				req.Header.Set("X-Forwarded-Proto", "https")
				req.Header.Set("Host", vhost.Host)
				req.Header.Set("X-Forwarded-For", req.RemoteAddr)
				req.Header.Set("X-Real-IP", req.RemoteAddr)
				req.Header.Set("X-Forwarded-Port", "443")
				req.Header.Set("X-Forwarded-SSL", "on")
				req.Header.Set("X-Forwarded-Ssl", "on")
			}

			//Header response
			proxy.ModifyResponse = func(resp *http.Response) error {
				resp.Header.Set("Server", "Setkaro")
				resp.Header.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
				resp.Header.Set("Alt-Svc", "h2=\":443\"; ma=2592000")
				resp.Header.Set("X-Forwarded-Proto", "https")
				resp.Header.Set("Content-Security-Policy", "upgrade-insecure-requests")
				return nil
			}

			//Allow to use cookie
			proxy.Transport = &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			}

			//Insecure skip verify
			proxy.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

			mux.HandleFunc(domain.Domain+"/", func(w http.ResponseWriter, r *http.Request) {
				//Copy cors header
				w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
				//Copy cookies header
				w.Header().Set("Cookie", r.Header.Get("Cookie"))
				//Copy credentials header
				w.Header().Set("Access-Control-Allow-Credentials", r.Header.Get("Credentials"))

				proxy.ServeHTTP(w, r)
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
		//Alow JSON header
		req.Header.Set("Content-Type", "application/json")
		req.URL.Scheme = vhost.Scheme
		req.URL.Host = vhost.Host
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Forwarded-Proto", "https")
		req.Header.Set("Host", vhost.Host)
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		req.Header.Set("X-Real-IP", req.RemoteAddr)
		req.Header.Set("X-Forwarded-Port", "443")
		req.Header.Set("X-Forwarded-SSL", "on")
		req.Header.Set("X-Forwarded-Ssl", "on")
	}

	//Header response
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Server", "Setkaro")
		return nil
	}

	//Allow to use cookie
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	mux.HandleFunc(hostname+"/", func(w http.ResponseWriter, r *http.Request) {
		//Allow method
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		proxy.ServeHTTP(w, r)
	})

	//Add default domain
	domains = append(domains, hostname)

	//Run on port 80
	go http.ListenAndServe(":80", mux)
	//Listen and serve
	autotls.Run(mux, domains...)
}
