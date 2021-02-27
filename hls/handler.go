package hls

import (
	"fmt"
	"net/http"
)

func (s *Server) videoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parameters := r.URL.Query()

		videoIDs, ok := parameters[videoIDKey]
		if !ok {
			s.writeInfoResponse(
				w,
				r,
				[]byte(fmt.Sprintf("Video ID url parameter ('vid') is missing")),
				http.StatusInternalServerError,
				map[string]string{
					contentType: contentTypeTextHTML,
				},
			)

			s.logger.Info().Msg(fmt.Sprintf("Request with video ID url parameter missing"))
			return
		}

		// Multiple video ids currently not supported.
		http.ServeFile(w, r, s.videoIndex(videoIDs[0]))
		w.Header().Set(contentType, contentTypeM3U8)

		s.logger.Info().Msg(fmt.Sprintf("Request for video '%s' got served", videoIDs[0]))
	})
}

func (s *Server) videoIndex(filename string) string {
	return fmt.Sprintf("content/video/%s.m3u8", filename)
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
