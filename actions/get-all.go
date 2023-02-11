package actions

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Domains struct {
	Domain string `json:"domain"`
	Proxy  string `json:"proxy"`
}

func GetAll() ([]Domains, error) {
	//Create context
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	//Get all proxies
	cursor, err := ProxyCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//Create array of proxies
	var proxies []bson.M

	//Loop through all proxies
	for cursor.Next(ctx) {
		var proxy bson.M
		err = cursor.Decode(&proxy)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		proxies = append(proxies, proxy)
	}

	//Close the cursor
	cursor.Close(ctx)

	//Create array of domains
	var domains []Domains

	//Loop through all proxies
	for _, proxy := range proxies {
		domainValue, ok := proxy["domain"]
		if !ok {
			fmt.Println("domain field not found in proxy")
			continue
		}
		domain, ok := domainValue.(string)
		if !ok {
			fmt.Printf("domain field is not a string, got %T\n", domainValue)
			continue
		}

		proxyValue, ok := proxy["proxy"]
		if !ok {
			fmt.Println("proxy field not found in proxy")
			continue
		}
		prox, ok := proxyValue.(string)
		if !ok {
			fmt.Printf("proxy field is not a string, got %T\n", proxyValue)
			continue
		}

		domains = append(domains, Domains{
			Domain: domain,
			Proxy:  prox,
		})
	}

	return domains, nil
}
