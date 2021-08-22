package handlers

import "net/http"

func ApiCallback(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
