package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	"golang.org/x/crypto/acme/autocert"
)

// type Domains struct {
// 	Domain string `json:"domain"`
// 	Proxy  string `json:"proxy"`
// 	SSL    bool   `json:"ssl"`
// }

// func main() {

// 	//Read JSON file
// 	jsonFile, err := os.ReadFile("dauqu.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Successfully Opened dauqu.json")

// 	var dauqu []Domains

// 	//Unmarshal JSON file
// 	err = json.Unmarshal(jsonFile, &dauqu)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	//Lpp through domains
// 	for _, domain := range dauqu {

// 		vhost, err := url.Parse(domain.Proxy)
// 		if err != nil {
// 			panic(err)
// 		}

// 		//Create proxy
// 		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
// 			Scheme: "http",
// 			Host:   domain.Proxy,
// 		})

// 		//Allow use https
// 		proxy.ModifyResponse = func(resp *http.Response) error {
// 			resp.Header.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
// 			return nil
// 		}

// 		//Set Header
// 		proxy.Director = func(req *http.Request) {
// 			req.Header.Set("X-Forwarded-Host", req.Host)
// 			req.Header.Set("X-Origin-Host", vhost.Host)
// 			req.Header.Set("X-Forwarded-Proto", "https")
// 			req.Header.Set("X-Forwarded-Port", "443")
// 			req.Header.Set("X-Forwarded-For", req.RemoteAddr)
// 			req.URL.Scheme = vhost.Scheme
// 			req.URL.Host = vhost.Host
// 		}

// 		// RunWithManagerAndTLSConfig support custom autocert manager and tls.Config

// 		//Create autocert manager
// 		m := autocert.Manager{
// 			Prompt:     autocert.AcceptTOS,
// 			HostPolicy: autocert.HostWhitelist(domain.Domain),
// 			Cache:      autocert.DirCache("certs"),
// 		}

// 		//Create TLS config
// 		tlsConfig := &tls.Config{
// 			GetCertificate: m.GetCertificate,
// 		}

// 		//Create server
// 		s := &http.Server{
// 			Addr:      ":443",
// 			TLSConfig: tlsConfig,
// 		}

// 		//Create handler

// 		//Create route

// 		http.HandleFunc(domain.Domain, handler(proxy))
// 	}

// 	//Listen on port 80
// 	http.ListenAndServe(":80", nil)
// }

// func handler(p *httputil.ReverseProxy) func(http.ResponseWriter,
// 	*http.Request) {

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		p.ServeHTTP(w, r)
// 	}
// }

func main() {

	mux := http.NewServeMux()

	vhost, err := url.Parse("http://localhost:44593")
	if err != nil {
		panic(err)
	}

	//Create proxy
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: vhost.Scheme,
		Host:   vhost.Host,
	})

	//Create route
	mux.HandleFunc("a.setkaro.com", handler(proxy))

	//Create autocert manager
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("a.setkaro.com" ),
		Cache:      autocert.DirCache("certs"),
	}

	//Create TLS config
	tlsConfig := &tls.Config{
		GetCertificate: m.GetCertificate,
	}

	//Create server
	s := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler:   mux,
	}

	//Listen on port 80
	go http.ListenAndServe(":80", m.HTTPHandler(nil))

	//Listen on port 443
	s.ListenAndServeTLS("", "")
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter,
	*http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}
