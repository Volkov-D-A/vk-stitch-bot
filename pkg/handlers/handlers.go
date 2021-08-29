package handlers

import (
	"log"
	"net/http"
)

//ApiCallback base handler for callback requests from VK api
func ApiCallback(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello, world!")); if err != nil {
		log.Print(err)
	}
}
