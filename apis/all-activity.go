package apis

import (
	"context"
	"dauqu-server/config"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ProxyCollection *mongo.Collection = config.GetCollection(config.DB, "counter")

func AllActivity(w http.ResponseWriter, r *http.Request) {

	//Create context
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	//Get all proxies and return them
	cursor, err := ProxyCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}

	//Get all proxies and return them
	var results []bson.M

	if err = cursor.All(ctx, &results); err != nil {
		fmt.Println(err)
	}

	//Return all JSOn
	type Counter struct {
		Id       string `json:"id"`
		Ip       string `json:"ip"`
		Hostname string `json:"hostname"`
		Port     string `json:"port"`
		Method   string `json:"method"`
		Time     string `json:"time"`
	}

	var counters []Counter

	for _, result := range results {
		counters = append(counters, Counter{
			Id:       result["_id"].(string),
			Ip:       result["ip"].(string),
			Hostname: result["hostname"].(string),
			Port:     result["port"].(string),
			Method:   result["method"].(string),
			Time:     result["time"].(string),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(counters)
}

func Analytics(w http.ResponseWriter, r *http.Request) {

	//Create context
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	//GEt length of all proxies
	total_requests, err := ProxyCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
	}

	//Get unique visitors
	unique_visitors, err := ProxyCollection.Distinct(ctx, "ip", bson.D{})
	if err != nil {
		fmt.Println(err)
	}

	//Get 24 hours
	hours_24, err := ProxyCollection.CountDocuments(ctx, bson.M{"time": bson.M{"$gt": time.Now().Add(-24 * time.Hour)}})
	if err != nil {
		fmt.Println(err)
	}

	//Return all JSOn
	type Total struct {
		TotalRequests  int `json:"total_requests"`
		UniqueVisitors int `json:"unique_visitors"`
		Last24Hours    int `json:"last_24_hours"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Total{
		TotalRequests:  int(total_requests),
		UniqueVisitors: len(unique_visitors),
		Last24Hours:    int(hours_24),
	})
}

// analytics by hostname
func AnalyticsByPort(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Port string `json:"port"`
	}

	var port Body
	err := json.NewDecoder(r.Body).Decode(&port)
	if err != nil {
		fmt.Println(err)
	}

	//Create context
	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	//Get all proxies and return them
	total_requests, err := ProxyCollection.CountDocuments(ctx, bson.M{"port": port.Port})
	if err != nil {
		fmt.Println(err)
	}

	//Get all proxies and return them
	unique_visitors, err := ProxyCollection.Distinct(ctx, "ip", bson.M{"port": port.Port})
	if err != nil {
		fmt.Println(err)
	}

	//Get all proxies and return them
	hours_24, err := ProxyCollection.CountDocuments(ctx, bson.M{"port": port.Port, "time": bson.M{"$gt": time.Now().Add(-24 * time.Hour)}})
	if err != nil {
		fmt.Println(err)
	}

	//Return all JSOn
	type Total struct {
		TotalRequests  int `json:"total_requests"`
		UniqueVisitors int `json:"unique_visitors"`
		Last24Hours    int `json:"last_24_hours"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Total{
		TotalRequests:  int(total_requests),
		UniqueVisitors: len(unique_visitors),
		Last24Hours:    int(hours_24),
	})
}
