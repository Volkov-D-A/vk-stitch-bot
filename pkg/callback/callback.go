package callback

import (
	"context"
	"net/http"
	"time"
)

//Server - base struct for callback server
type Server struct {
	callbackServer *http.Server
}

//Run the callbackServer instance
func (s *Server) Run(mux *http.ServeMux, port string) error {
	s.callbackServer = &http.Server{
		Addr:           ":" + port,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        mux,
	}
	return s.callbackServer.ListenAndServeTLS("/cert/server.crt", "/cert/server.key")
}

//Shutdown the callbackServer instance
func (s *Server) Shutdown(ctx context.Context) error {
	return s.callbackServer.Shutdown(ctx)
}
