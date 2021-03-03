package hls

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// videoPath constructs video file path for the handler.
func (s *Server) videoPath(filename string) string {
	return fmt.Sprintf("content/video/%s", filename)
}

// hlsVideoHandler defines a hls video handler.
func (s *Server) hlsVideoHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		videoID := chi.URLParam(r, "videoID")
		segmentID := chi.URLParam(r, "segmentID")

		if segmentID != "" {
			http.ServeFile(w, r, s.videoPath(segmentID))
			w.Header().Set(contentType, contentTypeTS)
			s.logger.Info().Msg(fmt.Sprintf("Request for video '%s' and segment '%s' got served", videoID, segmentID))
		} else {
			http.ServeFile(w, r, s.videoPath(videoID))
			w.Header().Set(contentType, contentTypeM3U8)
			s.logger.Info().Msg(fmt.Sprintf("Request for video '%s' got served", videoID))
		}
	})
}

// writeInfoResponse is a helper method for returning info http responses.
func (s *Server) writeInfoResponse(
	w http.ResponseWriter,
	r *http.Request,
	message []byte,
	status int,
	headers map[string]string,
) {
	for k, v := range headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(status)
	w.Write([]byte(message))
}
