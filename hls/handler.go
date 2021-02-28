package hls

import (
	"fmt"
	"net/http"
)

func (s *Server) videoPath() string {
	return fmt.Sprintf("content/video")
}

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
