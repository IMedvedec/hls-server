package hls

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
)

// Server defines an application server with its dependencies.
type Server struct {
	logger *zerolog.Logger
	server *http.Server
}

// New is a Server constructor.
func New(
	address string,
) *Server {
	var server Server

	mux := http.NewServeMux()
	// video file server.
	mux.Handle(server.videoPath(), http.StripPrefix(server.videoPath(), http.FileServer(http.Dir(server.videoPath()))))

	httpServer := http.Server{
		Addr:    address,
		Handler: mux,
	}

	consoleWriter := zerolog.NewConsoleWriter()
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	server.logger = &logger
	server.server = &httpServer
	return &server
}

// ListenAndServe acts as a http listen and serve wrapper.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Shutdown acts as a http shutdown wrapper.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
