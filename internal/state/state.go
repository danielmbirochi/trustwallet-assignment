package state

// KeyValueStorer is an interface for defining key-value store behavior.
type KeyValueStorer interface {
	Has(key string) (bool, error)
	Get(key string) ([][]byte, error)
	Put(key string, value [][]byte) error
	Delete(key string) error
}
