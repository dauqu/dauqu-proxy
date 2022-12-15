package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"

	// "github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

// Create JSON data
var jsonData = []byte(`[
	{
		"domain": "localhost/1",
		"proxy": "localhost:58441",
		"headers": ["X-Forwarded-For", "X-Forwarded-Proto", "X-Forwarded-Host", "X-Forwarded-Port"],
		"ssl": true
	},
	{
		"domain": "localhost/2",
		"proxy": "localhost:56710",
		"headers": ["X-Forwarded-For", "X-Forwarded-Proto", "X-Forwarded-Host", "X-Forwarded-Port"],
		"ssl": false
	}
]`)

func main() {

	//Create Details struct
	type Details struct {
		Domain  string   `json:"domain"`
		Proxy   string   `json:"proxy"`
		Headers []string `json:"headers"`
		SSL     bool     `json:"sll"`
	}

	//Bind JSON to Details struct
	var details []Details

	//Unmarshal JSON
	err := json.Unmarshal(jsonData, &details)
	if err != nil {
		log.Fatal(err)
	}

	//Print details
	fmt.Println(details)

	//Loop through details
	for _, i := range details {

		//Create Gin router
		router := gin.Default()

		//Create proxy
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   i.Proxy,
		})

		//Create proxy handler
		proxyHandler := func(c *gin.Context) {
			proxy.ServeHTTP(c.Writer, c.Request)
		}

		//Create proxy route
		router.Any("/*path", proxyHandler)

		//Add headers
		// for _, header := range details.Headers {
		// 	router.Use(func(c *gin.Context) {
		// 		c.Writer.Header().Set(header, c.Request.Header.Get(header))
		// 		c.Next()
		// 	})
		// }

		//Start server
		// if i.SSL {
		// 	autotls.Run(router, i.Domain)
		// } else {
		// 	router.Run(i.Domain)
		// }

		router.Run(":80")
	}
}
