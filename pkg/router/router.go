package router

import (
	"net/http"
	"projects/vk-stitch-bot/pkg/handlers"
)

func GetRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", handlers.ApiCallback)
	return mux
}