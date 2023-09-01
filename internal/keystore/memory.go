package keystore

import (
	"errors"

	"github.com/bencleary/uploader"
)

var _ uploader.KeyStoreService = (*InMemoryKeyStore)(nil)

type InMemoryKeyStore struct {
	storage map[string][]byte
}

func NewInMemoryKeyStore() *InMemoryKeyStore {
	return &InMemoryKeyStore{
		storage: make(map[string][]byte),
	}
}

func (k *InMemoryKeyStore) StoreKey(id string, key []byte) error {
	k.storage[id] = key
	return nil
}

func (k *InMemoryKeyStore) RetrieveKey(id string) ([]byte, error) {
	key, ok := k.storage[id]
	if !ok {
		return nil, errors.New("key not found")
	}
	return key, nil
}

func (k *InMemoryKeyStore) DeleteKey(id string) error {
	delete(k.storage, id)
	return nil
}
