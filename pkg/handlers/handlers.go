package handlers

import "net/http"

//ApiCallback base handler for callback requests from VK api
func ApiCallback(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
