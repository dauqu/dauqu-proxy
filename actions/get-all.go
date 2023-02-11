package actions

import (
	"context"
	"dauqu-server/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ProxyCollection *mongo.Collection = config.GetCollection(config.DB, "proxies")

// Return rows as JSON
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
		return nil, err
	}

	//Bind rows to array
	var domains []Domains

	//Loop through all rows
	for cursor.Next(ctx) {
		var domain Domains
		err = cursor.Decode(&domain)
		if err != nil {
			return nil, err
		}

		domains = append(domains, domain)
	}

	return domains, nil
}
