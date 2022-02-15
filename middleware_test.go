package cache_handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/StevenCyb/cache_handler/store"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	middlewareTestCounter := 0

	store := store.NewInMemoryStore(time.Second * 2)
	testMiddlewareHandler := NewMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			middlewareTestCounter++
			w.Write([]byte(strconv.Itoa(middlewareTestCounter)))
		}, store,
		UseMethodKey{},
		UseHeaderKey{Key: "Authorization"},
		UseQueryParamsKey{Key: "name"},
	)

	request(t, &testMiddlewareHandler, "GET", "/", http.Header{}, 1)
	request(t, &testMiddlewareHandler, "GET", "/", http.Header{}, 1)

	time.Sleep(time.Second * 3)
	request(t, &testMiddlewareHandler, "GET", "/", http.Header{}, 2)
	request(t, &testMiddlewareHandler, "GET", "/", http.Header{}, 2)
	request(t, &testMiddlewareHandler, "POST", "/", http.Header{}, 3)
	request(t, &testMiddlewareHandler, "POST", "/", http.Header{}, 3)
	request(t, &testMiddlewareHandler, "POST", "/a", http.Header{}, 4)
	request(t, &testMiddlewareHandler, "POST", "/a", http.Header{}, 4)
	h := http.Header{}
	h.Add("Authorization", "a")
	request(t, &testMiddlewareHandler, "POST", "/a", h, 5)
	request(t, &testMiddlewareHandler, "POST", "/a", h, 5)
	h = http.Header{}
	h.Add("Authorization", "b")
	request(t, &testMiddlewareHandler, "POST", "/a", h, 6)
	request(t, &testMiddlewareHandler, "POST", "/a", h, 6)
	request(t, &testMiddlewareHandler, "POST", "/a?name=steven", h, 7)
	request(t, &testMiddlewareHandler, "POST", "/a?name=steven", h, 7)
	request(t, &testMiddlewareHandler, "POST", "/a?name=mike", h, 8)
	request(t, &testMiddlewareHandler, "POST", "/a?name=mike", h, 8)

	time.Sleep(time.Second * 3)
	request(t, &testMiddlewareHandler, "GET", "/", http.Header{}, 9)
	request(t, &testMiddlewareHandler, "POST", "/", http.Header{}, 10)
	request(t, &testMiddlewareHandler, "POST", "/a", http.Header{}, 11)
	h = http.Header{}
	h.Add("Authorization", "a")
	request(t, &testMiddlewareHandler, "POST", "/a", h, 12)
	h = http.Header{}
	h.Add("Authorization", "b")
	request(t, &testMiddlewareHandler, "POST", "/a", h, 13)
	request(t, &testMiddlewareHandler, "POST", "/a?name=steven", h, 14)
	request(t, &testMiddlewareHandler, "POST", "/a?name=mike", h, 15)
}

func request(t *testing.T, testMiddlewareHandler *http.HandlerFunc, method, path string, header http.Header, expectedCounter int) {
	errorDetails := fmt.Sprintf("[%s] %s - %+v", method, path, header)
	req, err := http.NewRequest(method, path, nil)
	assert.NoError(t, err, errorDetails)
	req.Header = header

	hr := NewHttpRecorder(httptest.NewRecorder())
	testMiddlewareHandler.ServeHTTP(hr, req)
	assert.Equal(t, strconv.Itoa(expectedCounter), hr.Body.String(), errorDetails)
}
