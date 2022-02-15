package store

// Store represents a store
type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
}
