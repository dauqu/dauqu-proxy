package main

import (
	"crypto/tls"
	actions "dauqu-server/actions"
	database "dauqu-server/config"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Domains struct {
	Domain string `json:"domain"`
	Proxy  string `json:"proxy"`
}

func main() {
	mux := http.NewServeMux()

	//Get Hostname
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}

	//Get domain from database
	db := database.Connect()

	//Create table if not exists
	_, err = db.Query("CREATE TABLE IF NOT EXISTS proxies (domain VARCHAR(255) NOT NULL, proxy VARCHAR(255) NOT NULL, PRIMARY KEY (domain))")
	if err != nil {
		fmt.Println(err)
	}

	//Get all proxies
	rows, err := db.Query("SELECT * FROM proxies")
	if err != nil {
		fmt.Println(err)
	}

	//Create array of domains
	var dauqu []Domains

	//Loop through rows
	for rows.Next() {
		var domain Domains

		//Scan rows
		err = rows.Scan(&domain.Domain, &domain.Proxy)
		if err != nil {
			fmt.Println(err)
		}

		//Add domain to domains
		dauqu = append(dauqu, domain)
	}

	//Close database connection
	database.Close(db)

	//Loop through domains
	if len(dauqu) > 0 {
		for _, domain := range dauqu {
			//Add domain to domains
			// domains = append(domains, domain.Domain)

			vhost, err := url.Parse(domain.Proxy)
			if err != nil {
				fmt.Println(err)
				continue
			}

			//Create proxy
			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: vhost.Scheme,
				Host:   vhost.Host,
			})

			//Set Header
			proxy.Director = func(req *http.Request) {
				req.URL.Scheme = vhost.Scheme
				req.URL.Host = vhost.Host
				req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
				req.Header.Set("X-Forwarded-Proto", "https")
				req.Header.Set("X-Forwarded-For", req.RemoteAddr)
				req.Header.Set("X-Real-IP", req.RemoteAddr)
				req.Header.Set("X-Forwarded-Port", "443")
				req.Header.Set("X-Forwarded-SSL", "on")
			}

			//Header response
			proxy.ModifyResponse = func(resp *http.Response) error {
				resp.Header.Set("Server", "Setkaro")
				resp.Header.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
				resp.Header.Set("Alt-Svc", "h2=\":443\"; ma=2592000")
				resp.Header.Set("X-Forwarded-Proto", resp.Header.Get("X-Forwarded-Proto"))
				resp.Header.Set("Content-Security-Policy", "upgrade-insecure-requests")
				resp.Header.Set("Access-Control-Allow-Origin", resp.Header.Get("Access-Control-Allow-Origin"))
				resp.Header.Set("Access-Control-Allow-Credentials", resp.Header.Get("Access-Control-Allow-Credentials"))
				resp.Header.Set("Content-Type", resp.Header.Get("Content-Type"))
				return nil
			}

			mux.HandleFunc(domain.Domain+"/", func(w http.ResponseWriter, r *http.Request) {
				// //Check if proxy is responding or not
				// _, err := http.Get(domain.Proxy)
				// if err != nil {
				// 	//Return html error
				// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
				// 	w.WriteHeader(http.StatusServiceUnavailable)
				// 	//Show html file
				// 	http.ServeFile(w, r, "/var/dauqu/dauqu-proxy/index.html")
				// } else {
					go func() {
						actions.Counter(r)
					}()
					proxy.ServeHTTP(w, r)
				// }
			})
		}
	}

	vhost, err := url.Parse("http://localhost:9000")
	if err != nil {
		fmt.Println(err)
	}

	//Create proxy
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: vhost.Scheme,
		Host:   vhost.Host,
	})

	//Set Header
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = vhost.Scheme
		req.URL.Host = vhost.Host
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Forwarded-Proto", "https")
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		req.Header.Set("X-Real-IP", req.RemoteAddr)
		req.Header.Set("X-Forwarded-Port", "443")
		req.Header.Set("X-Forwarded-SSL", "on")
	}

	//Header response
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Server", "Setkaro")
		resp.Header.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		resp.Header.Set("Alt-Svc", "h2=\":443\"; ma=2592000")
		resp.Header.Set("X-Forwarded-Proto", resp.Header.Get("X-Forwarded-Proto"))
		resp.Header.Set("Content-Security-Policy", "upgrade-insecure-requests")
		resp.Header.Set("Access-Control-Allow-Origin", resp.Header.Get("Access-Control-Allow-Origin"))
		resp.Header.Set("Access-Control-Allow-Credentials", resp.Header.Get("Access-Control-Allow-Credentials"))
		resp.Header.Set("Content-Type", resp.Header.Get("Content-Type"))
		return nil
	}

	//Allow to use cookie
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	mux.HandleFunc(hostname+"/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusOK)
			return
		}

		//Check if proxy is responding or not
		_, err := http.Get("http://localhost:9000")
		if err != nil {
			//Return html error
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			//Show html file
			http.ServeFile(w, r, "/var/dauqu/dauqu-proxy/index.html")
		} else {
			proxy.ServeHTTP(w, r)
		}
	})

	//Cert Manager for auto generate ssl
	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("/var/dauqu/cert"),
	}
	//Create server
	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	//Listen and serve http and https
	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")
}
