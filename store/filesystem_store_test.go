package store

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFilesystemStore(t *testing.T) {
	store := NewFilesystem("./", 1*time.Second)
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

	filesToCheck := []string{}
	for _, fileIndex := range store.fileIndex {
		filesToCheck = append(filesToCheck, fileIndex.path)
	}

	time.Sleep(1 * time.Second)
	_, err = store.Get("dummy1")
	assert.Error(t, err)
	_, err = store.Get("dummy2")
	assert.Error(t, err)

	time.Sleep(1 * time.Second)
	assert.Equal(t, 0, len(store.fileIndex))
	for _, fileToCheck := range filesToCheck {
		_, err := os.Stat(fileToCheck)
		assert.Equal(t, true, errors.Is(err, os.ErrNotExist),
			"FileSystemStoreGC failed to delete %s", fileToCheck)
	}
}
