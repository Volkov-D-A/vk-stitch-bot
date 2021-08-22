package callback

import (
	"context"
	"net/http"
	"projects/vk-stitch-bot/pkg/router"
	"time"
)

type Server struct {
	callbackServer *http.Server
}

func (s *Server) Run(port string) error {
	s.callbackServer = &http.Server{
		Addr: ":" + port,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler: router.GetRouter(),
	}
	return s.callbackServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.callbackServer.Shutdown(ctx)
}