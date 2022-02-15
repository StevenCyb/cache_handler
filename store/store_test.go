package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreInterface(t *testing.T) {
	var filesystemStore Store = NewFilesystem("", 0)
	assert.NotNil(t, filesystemStore)

	var inMemoryStore Store = NewInMemoryStore(0)
	assert.NotNil(t, inMemoryStore)

	var redisStore Store = NewRedisStore("", 0, "", "", 0)
	assert.NotNil(t, redisStore)
}
