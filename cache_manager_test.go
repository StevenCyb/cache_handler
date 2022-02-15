package cache_handler

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheManagerOptionUsage(t *testing.T) {
	useKeyOptions := []Options{
		UseHeaderKey{Key: "h1"},
		UseHeaderKey{Key: "h2"},
		UseMethodKey{},
		UsePathKey{},
		UseQueryParamsKey{Key: "qp1"},
		UseQueryParamsKey{Key: "qp2"},
	}
	allowBypassOptions := []Options{
		AllowBypassHeader{Key: "Cache-Status", Value: "bypass"},
		AllowBypassMethod{Key: "post"},
		AllowBypassMethod{Key: "put"},
	}

	cm := cacheManager{}
	cm.useOptions(allowBypassOptions...)
	cm.useOptions(useKeyOptions...)

	assert.Len(t, cm.IncludeKeyOptions, len(useKeyOptions))
	for i, ko := range useKeyOptions {
		assert.Equal(t, ko, cm.IncludeKeyOptions[i])
	}

	assert.Len(t, cm.BypassOptions, len(allowBypassOptions))
	for i, abo := range allowBypassOptions {
		assert.Equal(t, abo, cm.BypassOptions[i])
	}
}

func TestCacheManagerKeyFromRequest(t *testing.T) {
	keyOptions := []Options{
		UseMethodKey{},
		UsePathKey{},
		UseHeaderKey{Key: "H1"},
		UseHeaderKey{Key: "H2"},
		UseQueryParamsKey{Key: "qp1"},
		UseQueryParamsKey{Key: "qp2"},
	}
	cm := cacheManager{}
	cm.useOptions(keyOptions...)

	r, err := http.NewRequest("GET", "https://not-exists.com/sub?qp1=qp1&qp2=qp2", nil)
	assert.NoError(t, err)
	r.Header.Add("Cache-Status", "bypass")
	r.Header.Add("h1", "h1")
	r.Header.Add("h2", "h2")

	h := sha256.New()
	h.Write([]byte("GET//sub/h1/h2/qp1/qp2"))
	assert.Equal(t, hex.EncodeToString(h.Sum(nil)), cm.keyFromRequest(r))
}

func TestCacheManagerAllowBypassHeader(t *testing.T) {
	cm := cacheManager{
		BypassOptions: []Options{AllowBypassHeader{Key: "Cache-Status", Value: "bypass"}},
	}

	r, err := http.NewRequest("GET", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	r.Header.Add("Cache-Status", "bypass")
	assert.True(t, cm.canBypass(r))

	r, err = http.NewRequest("GET", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	r.Header.Add("Cache-Status", "wrong")
	assert.False(t, cm.canBypass(r))

	r, err = http.NewRequest("GET", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.False(t, cm.canBypass(r))
}

func TestCacheManagerAllowBypassMethod(t *testing.T) {
	cm := cacheManager{
		BypassOptions: []Options{AllowBypassMethod{Key: "post"}},
	}

	r, err := http.NewRequest("get", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.False(t, cm.canBypass(r))

	r, err = http.NewRequest("put", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.False(t, cm.canBypass(r))

	r, err = http.NewRequest("POST", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.True(t, cm.canBypass(r))

	r, err = http.NewRequest("post", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.True(t, cm.canBypass(r))
}
