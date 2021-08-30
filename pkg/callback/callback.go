package callback

import (
	"context"
	"net/http"
	"time"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/router"
)

//Server - base struct for callback server
type Server struct {
	callbackServer *http.Server
}

//Run the callbackServer instance
func (s *Server) Run(port string) error {
	s.callbackServer = &http.Server{
		Addr:           ":" + port,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        router.GetRouter(),
	}
	return s.callbackServer.ListenAndServe()
}

//Shutdown the callbackServer instance
func (s *Server) Shutdown(ctx context.Context) error {
	return s.callbackServer.Shutdown(ctx)
}
