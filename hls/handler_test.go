package hls

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

const (
	testServerHost string = "localhost:9999"
)

func TestHlsVideoHandler(t *testing.T) {
	cases := []struct {
		name               string
		requestPath        string
		requestMethod      string
		responseStatusCode int
		responseBody       string
	}{
		{
			name:               "valid stream get",
			requestPath:        "/multimedia/video/ocean_test.m3u8/stream",
			requestMethod:      http.MethodGet,
			responseStatusCode: http.StatusOK,
			responseBody: "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:10\n#EXT-X-MEDIA-SEQUENCE:0\n" +
				"#EXTINF:10.000000,\nhttp://localhost:9999/multimedia/video/ocean_test.m3u8/stream/ocean0.ts\n" +
				"#EXTINF:8.480000,\n" +
				"http://localhost:9999/multimedia/video/ocean_test.m3u8/stream/ocean1.ts\n" +
				"#EXT-X-ENDLIST\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			request := httptest.NewRequest(c.requestMethod, c.requestPath, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("videoID", "ocean_test.m3u8")
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			response := httptest.NewRecorder()

			New(testServerHost).hlsVideoHandler().ServeHTTP(response, request)

			assert.Equal(t, c.responseStatusCode, response.Result().StatusCode, "status codes should match")
			assert.Equal(t, []string{contentTypeM3U8}, response.Result().Header[contentType], "content types should match")
			assert.Equal(t, c.responseBody, response.Body.String(), "body content should be the same")
		})
	}
}
