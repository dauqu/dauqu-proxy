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

	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	cursor, err := ProxyCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		fmt.Println(err)
	}

	var domains []Domains
	for _, result := range results {
		domains = append(domains, Domains{
			Domain: result["domain"].(string),
			Proxy:  result["proxy"].(string),
		})
	}

	return domains, nil
}
