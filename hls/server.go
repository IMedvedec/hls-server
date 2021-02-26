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
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    address,
		Handler: mux,
	}

	consoleWriter := zerolog.NewConsoleWriter()
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	return &Server{
		logger: &logger,
		server: &server,
	}
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
