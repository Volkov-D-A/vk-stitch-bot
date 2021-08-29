package router

import (
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/handlers"
	"net/http"
)

//GetRouter returns the the ServeMux with callback handlers
func GetRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", handlers.ApiCallback)
	return mux
}
