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
	rows, err := db.Query("SELECT * FROM counters")
	if err != nil {
		fmt.Println(err)
	}

	type Counter struct {
		Ip      string `json:"ip"`
		Hostname string `json:"hostname"`
		Method  string `json:"method"`
		Time    string `json:"time"`
	}

	//Create array of domains
	var dauqu []Counter

	if len(dauqu) > 0 {
		for rows.Next() {
			var ip string
			var hostname string
			var method string
			var time string

			err = rows.Scan(&ip, &hostname, &method, &time)
			if err != nil {
				fmt.Println(err)
			}

			dauqu = append(dauqu, Counter{Ip: ip, Hostname: hostname, Method: method, Time: time})
		}
	}

	defer database.Close(db)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dauqu)

}