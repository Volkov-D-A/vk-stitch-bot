package handlers

import (
	"net/http"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/logs"
)

//ApiCallback base handler for callback requests from VK api
func ApiCallback(w http.ResponseWriter, r *http.Request) {
	logger := logs.Get()
	logger.Info(r)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		logger.Error(err)
	}
}
