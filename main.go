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
	proxy := router.Group(details.Domain)
	{
		proxy.Use(func(c *gin.Context) {
			for _, header := range details.Headers {
				c.Request.Header.Set(header, c.Request.Header.Get(header))
			}
			c.Next()
		})
		
		proxy.Any("/*path", func(c *gin.Context) {
			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",
				Host:   details.Proxy,
			})
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	//Start server
	if details.SLL {
		autotls.Run(router, details.Domain)
	} else {
		router.Run(":80")
	}
}