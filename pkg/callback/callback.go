package callback

import (
	"context"
	"net/http"
	"time"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"
)

var (
	cfg = config.GetConfig()
)

//Server - base struct for callback server
type Server struct {
	callbackServer *http.Server
}

//Run the callbackServer instance
func (s *Server) Run(mux *http.ServeMux) error {
	s.callbackServer = &http.Server{
		Addr:           ":" + cfg.CallbackPort,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        mux,
	}
	return s.callbackServer.ListenAndServe()
}

//Shutdown the callbackServer instance
func (s *Server) Shutdown(ctx context.Context) error {
	return s.callbackServer.Shutdown(ctx)
}
