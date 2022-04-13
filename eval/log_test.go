package eval

import (
	"testing"

	"github.com/puyihua/meme-cache/internal/node/store"
)

func TestSyncWalThroughput(t *testing.T) {
	numKeys := 200
	keys, err1 := readKeysFromJson(numKeys)
	if err1 != nil {
		t.Error(err1)
	}

	err2 := startNodeAndPutKeys(keys, store.TypeSyncWAL)
	if err2 != nil {
		t.Error(err2)
	}

	ops := []string{"PUT"}
	evalNodeThroughput(t, numKeys, 10, keys, ops)
}

func TestAsyncWalThroughput(t *testing.T) {
	numKeys := 200
	keys, err1 := readKeysFromJson(numKeys)
	if err1 != nil {
		t.Error(err1)
	}

	err2 := startNodeAndPutKeys(keys, store.TypeAsyncWAL)
	if err2 != nil {
		t.Error(err2)
	}

	ops := []string{"PUT"}
	evalNodeThroughput(t, numKeys, 10, keys, ops)
}
