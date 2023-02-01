package apis

import (
	database "dauqu-server/config"
	"encoding/json"
	"fmt"
	"net/http"
)

func AllActivity(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()

	//Get all proxies and return them
	rows, err := db.Query("SELECT * FROM counter")
	if err != nil {
		fmt.Println(err)
	}

	//Check database have 0 rows

	type Counter struct {
		Id       int    `json:"id"`
		Ip       string `json:"ip"`
		Hostname string `json:"hostname"`
		Port     string `json:"port"`
		Method   string `json:"method"`
		Time     string `json:"time"`
	}

	//Create array of domains
	var dauqu []Counter

	for rows.Next() {
		var id int
		var ip string
		var hostname string
		var port string
		var method string
		var time string

		err = rows.Scan(&id, &ip, &hostname, &port, &method, &time)
		if err != nil {
			fmt.Println(err)
		}

		dauqu = append(dauqu, Counter{Id: id, Ip: ip, Hostname: hostname, Port: port, Method: method, Time: time})
	}

	defer database.Close(db)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dauqu)

}

func Analytics(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()

	var total_requests int
	//Get all proxies and return them
	err := db.QueryRow("SELECT COUNT(*) FROM counter").Scan(&total_requests)
	if err != nil {
		fmt.Println(err)
	}

	var unique_visitors int
	//Get all proxies and return them
	err = db.QueryRow("SELECT COUNT(DISTINCT ip) FROM counter").Scan(&unique_visitors)
	if err != nil {
		fmt.Println(err)
	}

	var hours_24 int
	//Get all proxies and return them
	err = db.QueryRow("SELECT COUNT(*) FROM counter WHERE time > DATE_SUB(NOW(), INTERVAL 24 HOUR)").Scan(&hours_24)
	if err != nil {
		fmt.Println(err)
	}

	//Return all JSOn
	type Total struct {
		TotalRequests  int `json:"total_requests"`
		UniqueVisitors int `json:"unique_visitors"`
		Last24Hours    int `json:"last_24_hours"`
	}

	defer database.Close(db)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Total{TotalRequests: total_requests, UniqueVisitors: unique_visitors, Last24Hours: hours_24})
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

	db := database.Connect()

	var total_requests int
	//Get all proxies and return them
	err = db.QueryRow("SELECT COUNT(*) FROM counter WHERE port = ?", port.Port).Scan(&total_requests)
	if err != nil {
		fmt.Println(err)
	}

	var unique_visitors int
	//Get all proxies and return them
	err = db.QueryRow("SELECT COUNT(DISTINCT ip) FROM counter WHERE port = ?", port.Port).Scan(&unique_visitors)
	if err != nil {
		fmt.Println(err)
	}

	var hours_24 int
	//Get all proxies and return them
	err = db.QueryRow("SELECT COUNT(*) FROM counter WHERE time > DATE_SUB(NOW(), INTERVAL 24 HOUR) AND port = ?", port.Port).Scan(&hours_24)
	if err != nil {
		fmt.Println(err)
	}

	//Return all JSOn
	type Total struct {
		TotalRequests  int `json:"total_requests"`
		UniqueVisitors int `json:"unique_visitors"`
		Last24Hours    int `json:"last_24_hours"`
	}

	defer database.Close(db)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Total{TotalRequests: total_requests, UniqueVisitors: unique_visitors, Last24Hours: hours_24})
}
