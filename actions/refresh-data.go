package actions

import (
	"os"
)

func RefreshData() error {
	//Get all proxies
	dauqu, err := GetAll()
	if err != nil {
		return err
	}

	//Remove JSOn file and create new one
	err = os.Remove("/var/dauqu/dauqu.json")
	if err != nil {
		return err
	}

	//Loop dauqu and append to JSON file
	for _, domain := range dauqu {
		//Append data in JSON file
		err = os.WriteFile("/var/dauqu/dauqu.json", []byte(domain.Domain), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
