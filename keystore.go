package uploader

type KeyStoreService interface {
	// StoreKey stores a key with the given ID.
	StoreKey(id string, key []byte) error

	// RetrieveKey retrieves a key by its ID.
	RetrieveKey(id string) ([]byte, error)

	// DeleteKey deletes a key by its ID.
	DeleteKey(id string) error
}
