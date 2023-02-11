package config

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoDb
func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:7388139606@localhost:27017"))
	if err != nil {
		fmt.Println(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connected to MongoDB!")
	}
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("dauqu-proxy").Collection(collectionName)
	return collection
}
