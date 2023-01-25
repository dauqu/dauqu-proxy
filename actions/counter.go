package actions

import (
	"net/http"
)

// Create function that can accept request and response
func Counter(r *http.Request) {
	//Get IP address
	ip := r.RemoteAddr
	hostname := r.Host
	method := r.Method

	//Create mysql table if not exist	
	
}
