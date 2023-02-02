package actions

import (
	database "dauqu-server/config"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var conn *websocket.Conn

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	clients[ws] = true
	conn = ws
}

// Create function that can accept request and response
func Counter(r *http.Request, port_number string) {

	//Connect to database
	db := database.Connect()

	//Get IP address
	ip_port := r.RemoteAddr
	ip := strings.Split(ip_port, ":")[0]
	hostname := r.Host
	method := r.Method

	//Create table if not exist
	db.Exec("CREATE TABLE IF NOT EXISTS `counter` (`id` int(11) NOT NULL AUTO_INCREMENT,`ip` varchar(255) NOT NULL,`hostname` varchar(255) NOT NULL,`port` varchar(255) NOT NULL,`method` varchar(255) NOT NULL,`time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")
	//Insert data to database
	db.Exec("INSERT INTO counter (ip, hostname, port, method) VALUES (?, ?, ?, ?)", ip, hostname, port_number, method)

	//Close database connection
	database.Close(db)

	//Send data to websocket
	if conn != nil {
		conn.WriteJSON(map[string]string{"ip": ip, "hostname": hostname, "method": method})
	}

}
