package router

import (
	"net/http"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/handlers"
)

//GetRouter returns the the ServeMux with callback handlers
func GetRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", handlers.ApiCallback)
	return mux
}
