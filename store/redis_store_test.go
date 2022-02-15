package store

import (
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestRedisStore(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	port, err := strconv.Atoi(mr.Port())
	assert.NoError(t, err)
	redisStore := NewRedisStore(mr.Addr(), port, "", "", 1*time.Second)

	err = redisStore.Set("dummy1", []byte("content1"))
	assert.NoError(t, err)
	err = redisStore.Set("dummy2", []byte("content2"))
	assert.NoError(t, err)

	data, err := redisStore.Get("dummy1")
	assert.NoError(t, err)
	assert.Equal(t, data, []byte("content1"))
	data, err = redisStore.Get("dummy2")
	assert.NoError(t, err)
	assert.Equal(t, data, []byte("content2"))

	// BUG cant test expired keys, miniredis seems to ignore
	// key expiration
}
