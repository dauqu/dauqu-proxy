package actions

import (
	"encoding/json"
	"os"
)

func RefreshData() error {
	//Get all proxies
	proxies, err := GetAll()
	if err != nil {
		return err
	}

	//Remove JSON file and create new one
	err = os.Remove("/var/dauqu/dauqu.json")
	if err != nil {
		return err
	}

	//Create array of domains
	var allDetails []Domains

	//Loop through proxies
	for _, proxy := range proxies {
		//Create domain
		domain := Domains{
			Domain: proxy.Domain,
			Proxy:  proxy.Proxy,
		}

		//Append domain to array
		allDetails = append(allDetails, domain)
	}

	//Encode the array of domains to a JSON string
	jsonData, err := json.Marshal(allDetails)
	if err != nil {
		return err
	}

	//Write to file
	err = os.WriteFile("/var/dauqu/dauqu.json", jsonData, os.FileMode(0644))
	if err != nil {
		return err
	}
	return nil
}
