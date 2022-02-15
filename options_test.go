package cache_handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsePathKey(t *testing.T) {
	r, err := http.NewRequest("POST", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)

	uk := UsePathKey{}
	assert.Equal(t, "/sub", uk.ExtractString(r))
}

func TestUseMethodKey(t *testing.T) {
	r, err := http.NewRequest("POST", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)

	uk := UseMethodKey{}
	assert.Equal(t, "POST", uk.ExtractString(r))
}

func TestUseQueryParamsKey(t *testing.T) {
	r, err := http.NewRequest("POST", "https://not-exists.com/sub?name=abc&role=a&role=b", nil)
	assert.NoError(t, err)

	uk := UseQueryParamsKey{Key: "name"}
	assert.Equal(t, "abc", uk.ExtractString(r))
	uk = UseQueryParamsKey{Key: "role"}
	assert.Equal(t, "a/b", uk.ExtractString(r))
	uk = UseQueryParamsKey{Key: "not_exists"}
	assert.Equal(t, "", uk.ExtractString(r))
}

func TestUseHeaderKey(t *testing.T) {
	r, err := http.NewRequest("POST", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	r.Header.Add("Cache-Status", "bypass")

	uk := UseHeaderKey{Key: "Cache-Status"}
	assert.Equal(t, "bypass", uk.ExtractString(r))
	uk = UseHeaderKey{Key: "not_exists"}
	assert.Equal(t, "", uk.ExtractString(r))
}

func TestAllowBypassHeader(t *testing.T) {
	abh := AllowBypassHeader{Key: "Cache-Status", Value: "bypass"}

	r, err := http.NewRequest("POST", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.False(t, abh.ExtractBool(r))

	r, err = http.NewRequest("POST", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	r.Header.Add("Cache-Status", "wrong")
	assert.False(t, abh.ExtractBool(r))

	r, err = http.NewRequest("POST", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	r.Header.Add("Cache-Status", "bypass")
	assert.True(t, abh.ExtractBool(r))
}

func TestAllowBypassMethod(t *testing.T) {
	abm := AllowBypassMethod{Key: "post"}

	r, err := http.NewRequest("get", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.False(t, abm.ExtractBool(r))

	r, err = http.NewRequest("post", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.True(t, abm.ExtractBool(r))

	r, err = http.NewRequest("POST", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.True(t, abm.ExtractBool(r))

	r, err = http.NewRequest("post", "https://not-exists.com/sub", nil)
	assert.NoError(t, err)
	assert.True(t, abm.ExtractBool(r))
}
