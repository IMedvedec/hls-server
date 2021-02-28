package hls

import (
	"fmt"
	"net/http"
)

// videoPath is a helper method for defining video file server paths.
func (s *Server) videoPath() string {
	return fmt.Sprintf("content/video")
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
