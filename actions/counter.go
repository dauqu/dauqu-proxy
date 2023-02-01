package actions

import (
	database "dauqu-server/config"
	// "fmt"
	"net/http"
	"strings"
	// "github.com/gorilla/websocket"
)

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// var clients = make(map[*websocket.Conn]bool)
// var conn *websocket.Conn

// func WebSocket(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// register client
// 	clients[ws] = true
// 	conn = ws // Set the global connection
// }

// Create function that can accept request and response
func Counter(r *http.Request) {
	//Connect to database
	db := database.Connect()

	//Get IP address
	ip_port := r.RemoteAddr
	ip := strings.Split(ip_port, ":")[0]
	hostname := r.Host
	method := r.Method

	//Create table if not exist
	db.Exec("CREATE TABLE IF NOT EXISTS `counter` (`id` int(11) NOT NULL AUTO_INCREMENT,`ip` varchar(255) NOT NULL,`hostname` varchar(255) NOT NULL,`method` varchar(255) NOT NULL,`time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")
	//Insert data to database
	db.Exec("INSERT INTO `counter` (`ip`, `hostname`, `method`) VALUES (?, ?, ?)", ip, hostname, method)

	//Close database connection
	database.Close(db)

	//Send data to websocket
	// if conn != nil {
	// 	conn.WriteJSON(ip)
	// }

}