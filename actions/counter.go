package actions

import (
	database "dauqu-server/config"
	"net/http"
)

// Create function that can accept request and response
func Counter(r *http.Request) {

	db := database.Connect()

	//Get IP address
	ip := r.RemoteAddr
	hostname := r.Host
	method := r.Method

	//Create table if not exist
	db.Exec("CREATE TABLE IF NOT EXISTS `counter` (`id` int(11) NOT NULL AUTO_INCREMENT,`ip` varchar(255) NOT NULL,`hostname` varchar(255) NOT NULL,`method` varchar(255) NOT NULL,`time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")
	//Insert data to database
	db.Exec("INSERT INTO `counter` (`ip`, `hostname`, `method`) VALUES (?, ?, ?)", ip, hostname, method)

	//Close database connection
	database.Close(db)
}
