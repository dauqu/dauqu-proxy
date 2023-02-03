package actions

import (
	database "dauqu-server/config"
	"fmt"
)

type Domains struct {
	Domain string `json:"domain"`
	Proxy  string `json:"proxy"`
}

func GetAll() ([]Domains, error) {
	db := database.Connect()

	//Check database connection is not nil
	if db != nil {
		CreateTable()
	} else {
		fmt.Println("Failed to connect to database")
	}

	//Get all proxies and return them
	rows, err := db.Query("SELECT * FROM proxies")
	if err != nil {
		fmt.Println(err)
	}

	//Create array of domains
	var dauqu []Domains

	for rows.Next() {
		var domain string
		var proxy string

		err = rows.Scan(&domain, &proxy)
		if err != nil {
			fmt.Println(err)
		}

		dauqu = append(dauqu, Domains{Domain: domain, Proxy: proxy})
	}

	defer database.Close(db)

	//Append data in JSON file 

	return dauqu, nil
}
