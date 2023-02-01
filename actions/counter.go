package actions

import (
	"context"
	database "dauqu-server/config"
	"fmt"
	"net/http"
	"strings"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Create function that can accept request and response
func Counter(r *http.Request) {

	// Get the websocket connection from the request context
	conn, _ := r.Context().Value("conn").(*websocket.Conn)

	// Send a message to the Ws function
	err := conn.WriteJSON("counter hit")
	if err != nil {
		fmt.Println(err)
	}

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
}

func Ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	// Create a context with the connection
	ctx := context.WithValue(r.Context(), "conn", conn)

	// Call the counter function
	Counter(r.WithContext(ctx))



}
