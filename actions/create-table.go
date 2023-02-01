package actions

import (
	database "dauqu-server/config"
	"fmt"
)

func CreateTable() {
	db := database.Connect()

	_, err := db.Query("CREATE TABLE IF NOT EXISTS proxies (domain VARCHAR(255) NOT NULL, proxy VARCHAR(255) NOT NULL, port VARCHAR(10) NOT NULL, PRIMARY KEY (domain))")
	if err != nil {
		fmt.Println(err)
	}

	defer database.Close(db)
}
