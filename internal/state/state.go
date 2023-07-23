package state

// KeyValueStorer is an interface for defining key-value store behavior.
type KeyValueStorer interface {
	// Has retrieves if a key is present in the key-value store. It
	// should return an error if datastore is not initialized.
	Has(key string) (bool, error)

	// Get retrieves the given key if it's present in the key-value store.
	// It will return an error if datastore is not initialized or the key
	// is not found.
	Get(key string) ([][]byte, error)

	// Put inserts the given value into the key-value store. It expect a
	// slice of slice of byte to enable adding multiple items to the key-value
	// list. To add just one item to the value list just pass one item into the
	// slice - []byte{item}. It will return an error if datastore is not initialized.
	Put(key string, value [][]byte) error

	// Delete removes the given key from the key-value store. It will return
	// an error if datastore is not initialized.
	Delete(key string) error
}
