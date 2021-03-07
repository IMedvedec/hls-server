package hls

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

const (
	testServerHost string = "localhost:9999"
)

func TestHlsVideoHandler(t *testing.T) {
	cases := []struct {
		name string

		// Request parameters.
		requestPath      string
		requestMethod    string
		requestVideoID   string
		requestSegmentID string
		// Response parameters.
		responseStatusCode int
		responseHeaders    []string
		responseBody       []byte
	}{
		{
			name:               "valid stream get",
			requestPath:        "/multimedia/video/ocean_test.m3u8/stream",
			requestVideoID:     "ocean_test.m3u8",
			requestSegmentID:   "",
			requestMethod:      http.MethodGet,
			responseStatusCode: http.StatusOK,
			responseHeaders:    []string{contentTypeM3U8},
			responseBody: func() []byte {
				file, err := os.OpenFile(videoPath("ocean_test.m3u8"), os.O_RDONLY, 0666)
				if err != nil {
					panic(err)
				}

				content, err := ioutil.ReadAll(file)
				if err != nil {
					panic(err)
				}

				return content
			}(),
		},
		{
			name:               "valid segment get",
			requestPath:        "/multimedia/video/ocean_test.m3u8/stream/ocean_test0.ts",
			requestVideoID:     "ocean_test.m3u8",
			requestSegmentID:   "ocean_test0.ts",
			requestMethod:      http.MethodGet,
			responseStatusCode: http.StatusOK,
			responseHeaders:    []string{contentTypeTS},
			responseBody: func() []byte {
				file, err := os.OpenFile(videoPath("ocean_test0.ts"), os.O_RDONLY, 0666)
				if err != nil {
					panic(err)
				}

				content, err := ioutil.ReadAll(file)
				if err != nil {
					panic(err)
				}

				return content
			}(),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			request := httptest.NewRequest(c.requestMethod, c.requestPath, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("videoID", c.requestVideoID)
			rctx.URLParams.Add("segmentID", c.requestSegmentID)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			response := httptest.NewRecorder()

			New(testServerHost).hlsVideoHandler().ServeHTTP(response, request)

			assert.Equal(t, c.responseStatusCode, response.Result().StatusCode, "status codes should match")
			assert.Equal(t, c.responseHeaders, response.Result().Header[contentType], "content types should match")
			assert.Equal(t, c.responseBody, response.Body.Bytes(), "body content should be the same")
		})
	}
}
