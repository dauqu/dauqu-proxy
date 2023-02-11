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
	}

	//Create array of proxies
	var proxies []bson.M

	//Loop through all proxies
	for cursor.Next(ctx) {
		var proxy bson.M
		err = cursor.Decode(&proxy)
		if err != nil {
			fmt.Println(err)
		}

		proxies = append(proxies, proxy)
	}

	//Create array of domains
	var domains []Domains

	//Loop through all proxies
	for _, proxy := range proxies {
		domains = append(domains, Domains{
			Domain: proxy["domain"].(string),
			Proxy:  proxy["proxy"].(string),
		})
	}

	return domains, nil
}
