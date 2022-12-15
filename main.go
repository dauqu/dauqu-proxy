package main

import (
	"encoding/json"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

//Create JSON data
var jsonData = []byte(`{
	"domain": "d.setkaro.com",
	"proxy": "http://localhost:5556",
	"headers": ["X-Forwarded-For", "X-Forwarded-Proto", "X-Forwarded-Host", "X-Forwarded-Port"],
	"sll": true
}`)


func main() {

	//Create Details struct
	type Details struct {
		Domain string `json:"domain"`
		Proxy  string `json:"proxy"`
		Headers []string `json:"headers"`
		SLL	bool `json:"sll"`
	}

	//Bind JSON to Details struct
	var details Details
	err := json.Unmarshal(jsonData, &details)
	if err != nil {
		log.Fatal(err)
	}


	//Create Gin router
	router := gin.Default()

	//Create proxy
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   details.Proxy,
	})

	//Create proxy handler
	proxyHandler := func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}

	//Create proxy route
	router.Any("/*path", proxyHandler)

	//Add headers
	for _, header := range details.Headers {
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set(header, c.Request.Header.Get(header))
			c.Next()
		})
	}

	//Start server
	if details.SLL {
		autotls.Run(router, details.Domain)
	} else {
		router.Run(":80")
	}
}