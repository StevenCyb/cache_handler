package cache_handler

import (
	"net/http"
	"strings"
)

// Options represent a option for the middleware
type Options interface {
	ExtractString(r *http.Request) string
	ExtractBool(r *http.Request) bool
}

// UsePathKey tells the cache to use path as part of the key
type UsePathKey struct{}

// ExtractString key from request and return the string
func (opt UsePathKey) ExtractString(r *http.Request) string {
	return r.URL.Path
}

// ExtractBool does nothing but return false
// (function required to match the interface)
func (opt UsePathKey) ExtractBool(r *http.Request) bool { return false }

// UseMethodKey tells the cache to use the http method as part of the key
type UseMethodKey struct{}

// ExtractString key from request and return the string
func (opt UseMethodKey) ExtractString(r *http.Request) string {
	return r.Method
}

// ExtractBool does nothing but return false
// (function required to match the interface)
func (opt UseMethodKey) ExtractBool(r *http.Request) bool { return false }

// UseQueryParamsKey tells the cache to use query parameter value for given key
// as part of the key
type UseQueryParamsKey struct{ Key string }

// ExtractString key from request and return the string
func (opt UseQueryParamsKey) ExtractString(r *http.Request) string {
	params := r.URL.Query()[opt.Key]
	return strings.Join(params, "/")
}

// ExtractBool does nothing but return false
// (function required to match the interface)
func (opt UseQueryParamsKey) ExtractBool(r *http.Request) bool { return false }

// UseQueryParamsKey tells the cache to........
type UseHeaderKey struct{ Key string }

// ExtractString key from request and return the string
func (opt UseHeaderKey) ExtractString(r *http.Request) string {
	params := r.Header[opt.Key]
	return strings.Join(params, "/")
}

// ExtractBool does nothing but return false
// (function required to match the interface)
func (opt UseHeaderKey) ExtractBool(r *http.Request) bool { return false }

// AllowBypassHeader enable a client to set the `Cache-Status`
// header to `bypass` so he will not get cached data
type AllowBypassHeader struct{ Key, Value string }

// ExtractString does nothing but return empty string
// (function required to match the interface)
func (opt AllowBypassHeader) ExtractString(r *http.Request) string { return "" }

// ExtractBool check if request can bypass depending on request
func (opt AllowBypassHeader) ExtractBool(r *http.Request) bool {
	if val, ok := r.Header[opt.Key]; ok && val[0] == opt.Value {
		for _, v := range val {
			if v == opt.Value {
				return true
			}
		}
	}

	return false
}

// AllowBypassMethod defined methods that automatically bypass the caching e.g. post or put
type AllowBypassMethod struct{ Key string }

// ExtractString does nothing but return empty string
// (function required to match the interface)
func (opt AllowBypassMethod) ExtractString(r *http.Request) string { return "" }

// ExtractBool check if request can bypass depending on request
func (opt AllowBypassMethod) ExtractBool(r *http.Request) bool {
	return opt.Key == strings.ToLower(r.Method)
}
