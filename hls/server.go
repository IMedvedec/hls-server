package hls

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	router := chi.NewRouter()
	// video file server.
	router.Route("/fileserver/content", func(r chi.Router) {
		r.Get("/video/{videoID}", http.StripPrefix("/fileserver/content/video/", http.FileServer(http.Dir("content/video"))).ServeHTTP)
	})
	// REST endpoints.
	router.Route("/multimedia/video", func(r chi.Router) {
		r.Get("/{videoID}/stream", server.hlsVideoHandler())
		r.Get("/{videoID}/stream/{segmentID}", server.hlsVideoHandler())
	})

	httpServer := http.Server{
		Addr:    address,
		Handler: router,
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
