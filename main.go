package main

import (
	"crypto/tls"
	actions "dauqu-server/actions"
	"dauqu-server/apis"
	database "dauqu-server/config"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Domains struct {
	Domain string `json:"domain"`
	Proxy  string `json:"proxy"`
}

func main() {
	mux := http.NewServeMux()

	//MongoDb
	database.ConnectDB()

	//Get Hostname
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}

	dauqu, err := actions.GetAll()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(dauqu)

	actions.RefreshData()

	// 	//Read JSON file
	jsonFile, err := os.ReadFile("/var/dauqu/dauqu.json")
	if err != nil {
		//If file not found create new file
		if os.IsNotExist(err) {
			_, err = os.Create("/var/dauqu/dauqu.json")
			if err != nil {
				fmt.Println(err)
			}

			//Write to file
			err = os.WriteFile("/var/dauqu/dauqu.json", []byte("[]"), 0644)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}

	fmt.Println("Successfully Opened dauqu.json")

	// var dauqu []Domains

	// 	//Unmarshal JSON file
	err = json.Unmarshal(jsonFile, &dauqu)
	if err != nil {
		fmt.Println(err)
	}

	//Loop through domains
	if len(dauqu) > 0 {
		for _, domain := range dauqu {
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

			var host_proxy = vhost.Host
			//Splite left side of : to get host
			port := strings.Replace(host_proxy, "localhost:", "", 1)

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
				//Set content type all text jpg png mp4 mp3 etc 
				resp.Header.Set("Content-Type", "*/*")
				return nil
			}

			mux.HandleFunc(domain.Domain+"/", func(w http.ResponseWriter, r *http.Request) {
				go func() {
					actions.Counter(r, port)
				}()
				proxy.ServeHTTP(w, r)
			})
		}
	}

	//API function
	mux.HandleFunc(hostname+"/dp/all-activity", apis.AllActivity)
	mux.HandleFunc(hostname+"/dp/analytics", apis.Analytics)
	mux.HandleFunc(hostname+"/dp/analytics-by-hostname", apis.AnalyticsByPort)
	mux.HandleFunc(hostname+"/dp/ws", actions.WsHandler)

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
			http.ServeFile(w, r, "/var/dauqu/server/index.html")
		} else {
			proxy.ServeHTTP(w, r)
		}
	})

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("/var/dauqu/cert"),
	}

	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")
}
