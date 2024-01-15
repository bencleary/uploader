package keystore_test

import (
	"testing"

	"github.com/bencleary/uploader/internal/keystore"
)

func TestNewMemoryKeyStore(t *testing.T) {
	ks := keystore.NewInMemoryKeyStore()
	if ks == nil {
		t.Fatal("expected keystore")
	}
}

func TestMemoryKeyStoreStoreKey(t *testing.T) {
	ks := keystore.NewInMemoryKeyStore()
	err := ks.StoreKey("key", []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestMemoryKeyStoreRetrieveKey(t *testing.T) {
	ks := keystore.NewInMemoryKeyStore()
	err := ks.StoreKey("key", []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	key, err := ks.RetrieveKey("key")
	if err != nil {
		t.Fatal("expected error")
	}
	if string(key) == "" {
		t.Fatal("expected key")
	}
}

func TestMemoryKeyStoreDeleteKey(t *testing.T) {
	ks := keystore.NewInMemoryKeyStore()
	err := ks.StoreKey("key", []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	res := ks.DeleteKey("key")
	if res != nil {
		t.Fatal("expected error")
	}
}
