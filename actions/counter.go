package actions

import (
	"context"
	"dauqu-server/config"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

var CounterCollection *mongo.Collection = config.GetCollection(config.DB, "counter")

// Create function that can accept request and response
func Counter(r *http.Request, port_number string) {

	//Get IP address
	ip_port := r.RemoteAddr
	ip := strings.Split(ip_port, ":")[0]
	hostname := r.Host
	method := r.Method

	//Insert data to database
	// db.Exec("INSERT INTO counter (ip, hostname, port, method) VALUES (?, ?, ?, ?)", ip, hostname, port_number, method)

	//Create context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Insert data to database
	_, err := CounterCollection.InsertOne(ctx, bson.M{"ip": ip, "hostname": hostname, "port": port_number, "method": method})
	if err != nil {
		fmt.Println(err)
	}

	//Send data to websocket
	if conn != nil {
		conn.WriteJSON(map[string]string{"ip": ip, "hostname": hostname, "method": method})
	}

}
