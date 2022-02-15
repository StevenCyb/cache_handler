package cache_handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpRecorder(t *testing.T) {
	bodyBytes := []byte("Sample body content we want to see.")
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write(bodyBytes)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	hr := NewHttpRecorder(httptest.NewRecorder())
	handler.ServeHTTP(hr, req)

	assert.Equal(t, bodyBytes, hr.Body.Bytes())
}
