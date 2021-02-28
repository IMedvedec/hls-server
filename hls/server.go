package hls

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
)

type Server struct {
	logger *zerolog.Logger
	server *http.Server
}

func New(
	address string,
) *Server {
	var server Server

	mux := http.NewServeMux()
	// video file server.
	mux.Handle("/content/video/", http.StripPrefix("/content/video/", http.FileServer(http.Dir(server.videoPath()))))

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

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
