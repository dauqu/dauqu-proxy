package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Domains struct {
	Domain  string   `json:"domain"`
	Proxy   string   `json:"proxy"`
	SSL     bool     `json:"ssl"`
}

func main() {

	//Read JSON file
	jsonFile, err := os.ReadFile("dauqu.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Opened dauqu.json")

	var dauqu []Domains

	//Unmarshal JSON file
	err = json.Unmarshal(jsonFile, &dauqu)
	if err != nil {
		log.Fatal(err)
	}

	//Loop through the JSON file
	for i := 0; i < len(dauqu); i++ {
		fmt.Println("Domain:", dauqu[i].Domain)
		fmt.Println("Proxy:", dauqu[i].Proxy)
		fmt.Println("SSL:", dauqu[i].SSL)
	}
}
