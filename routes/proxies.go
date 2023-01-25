package routes

import (
	database "dauqu-server/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

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

	//GO routine
	db := database.Connect()

	//Create table if not exists
	_, err = db.Query("CREATE TABLE IF NOT EXISTS proxies (domain VARCHAR(255) NOT NULL, proxy VARCHAR(255) NOT NULL, PRIMARY KEY (domain))")
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

	//Check if domain already exists
	var domain string
	err = db.QueryRow("SELECT domain FROM proxies WHERE domain = ?", bodyData.Domain).Scan(&domain)
	if err != nil {
		fmt.Println(err)
		//Continue
	}
	if domain == bodyData.Domain {
		//Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Domain already exists",
			"status":  http.StatusOK,
		})
		return
	}

	//Insert into database
	_, err = db.Query("INSERT INTO proxies (domain, proxy) VALUES (?, ?)", bodyData.Domain, bodyData.Proxy)
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

	//Close database connection
	database.Close(db)

	//Return response
	w.Header().Set("Content-Type", "application/json")
	//Return JSON response with error and message and status
	data := map[string]interface{}{
		"message": "Proxy added successfully",
		"status":  http.StatusOK,
	}
	json.NewEncoder(w).Encode(data)
}

// Get all proxies
func GetProxies(w http.ResponseWriter, r *http.Request) {

	//GO routine
	db := database.Connect()

	//Get all proxies
	rows, err := db.Query("SELECT * FROM proxies")
	if err != nil {
		fmt.Println(err)
	}

	//Return rows as JSON
	type Domains struct {
		Domain string `json:"domain"`
		Proxy  string `json:"proxy"`
	}

	var dauqu []Domains

	for rows.Next() {
		var domain Domains
		err = rows.Scan(&domain.Domain, &domain.Proxy)
		if err != nil {
			fmt.Println(err)
		}
		dauqu = append(dauqu, domain)
	}

	//Close database connection
	database.Close(db)

	//Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dauqu)
}

// Delete Proxy
func DeleteProxy(w http.ResponseWriter, r *http.Request) {
	
	type Body struct {
		Domain string `json:"domain"`
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

	//GO routine
	db := database.Connect()

	//Delete from database
	_, err = db.Query("DELETE FROM proxies WHERE domain = ?", bodyData.Domain)
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

	//Close database connection
	database.Close(db)

	//Return response
	w.Header().Set("Content-Type", "application/json")
	//Return JSON response with error and message and status
	data := map[string]interface{}{
		"message": "Proxy deleted successfully",
		"status":  http.StatusOK,
	}
	json.NewEncoder(w).Encode(data)
}