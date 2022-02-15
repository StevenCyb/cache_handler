package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryStore(t *testing.T) {
	store := NewInMemoryStore(1 * time.Second)
	err := store.Set("dummy1", []byte("content1"))
	assert.NoError(t, err)
	err = store.Set("dummy2", []byte("content2"))
	assert.NoError(t, err)

	data, err := store.Get("dummy1")
	assert.NoError(t, err)
	assert.Equal(t, data, []byte("content1"))
	data, err = store.Get("dummy2")
	assert.NoError(t, err)
	assert.Equal(t, data, []byte("content2"))

	time.Sleep(1 * time.Second)
	_, err = store.Get("dummy1")
	assert.Error(t, err)
	_, err = store.Get("dummy2")
	assert.Error(t, err)

	time.Sleep(1 * time.Second)
	assert.Equal(t, 0, len(store.data))
}
