package routes

import (
	"context"
	"dauqu-server/config"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"time"
)

var ProxyCollection *mongo.Collection = config.GetCollection(config.DB, "proxies")

// Add Proxy to the server
func AddProxy(w http.ResponseWriter, r *http.Request) {

	type Body struct {
		Domain string `json:"domain"`
		Proxy  string `json:"proxy"`
	}

	//Read Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err,
			"status":  http.StatusOK,
		})
		return
	}

	//Unmarshal Body
	var bodyData Body
	err = json.Unmarshal(body, &bodyData)
	if err != nil {
		//Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err,
			"status":  http.StatusOK,
		})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	//Check if proxy already exists
	result := ProxyCollection.FindOne(context.TODO(), bson.M{"domain": bodyData.Domain})
	if result.Err() == nil {
		//Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Proxy already exists",
			"status":  http.StatusOK,
		})
		return
	}

	//Insert into database
	resp, err := ProxyCollection.InsertOne(ctx, bson.M{
		"domain": bodyData.Domain,
		"proxy":  bodyData.Proxy,
	})
	if err != nil {
		//Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err,
			"status":  http.StatusOK,
		})
		return
	}

	//Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Proxy added successfully",
		"status":  http.StatusOK,
		"resp":    resp,
	})
}

// Get all proxies
func GetProxies(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	//Find all proxies
	cursor, err := ProxyCollection.Find(ctx, bson.M{})
	if err != nil {
		//Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err,
			"status":  http.StatusOK,
		})
		return
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

	//Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proxies)
}

// Delete Proxy
func DeleteProxy(w http.ResponseWriter, r *http.Request) {

	type Body struct {
		Domain string `json:"domain"`
		Proxy  string `json:"proxy"`
	}

	//Read Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err,
			"status":  http.StatusOK,
		})
		return
	}

	//Unmarshal Body
	var bodyData Body
	err = json.Unmarshal(body, &bodyData)
	if err != nil {
		//Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err,
			"status":  http.StatusOK,
		})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 600*time.Second)

	//Delete Proxy
	resp, err := ProxyCollection.DeleteOne(ctx, bson.M{
		"domain": bodyData.Domain,
		"proxy":  bodyData.Proxy,
	})
	if err != nil {
		//Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err,
			"status":  http.StatusOK,
		})
		return
	}

	//Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Proxy deleted successfully",
		"status":  http.StatusOK,
		"resp":    resp,
	})
}
