package cache_handler

import (
	"bytes"
	"net/http"
)

// HttpRecorder is a custom response writer that
// records the body
type HttpRecorder struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

// NewHttpRecorder create a new NewHttpRecorder with given ResponseWriter
func NewHttpRecorder(responseWriter http.ResponseWriter) *HttpRecorder {
	return &HttpRecorder{
		ResponseWriter: responseWriter,
		Body:           &bytes.Buffer{},
	}
}

// Write byte data is written to rw.Body, if not nil.
func (hr *HttpRecorder) Write(buf []byte) (int, error) {
	if hr.Body != nil {
		hr.Body.Write(buf)
		return hr.ResponseWriter.Write(buf)
	}

	return len(buf), nil
}

// WriteString string data is written to rw.Body, if not nil.
func (hr *HttpRecorder) WriteString(str string) (int, error) {
	if hr.Body != nil {
		hr.Body.WriteString(str)
		return hr.ResponseWriter.Write([]byte(str))
	}

	return len(str), nil
}
