package handler

import (
	"strings"

	"github.com/taymour/elysiandb/internal/storage"
	"github.com/taymour/elysiandb/internal/transport/tcp/parsing"
)

func HandleMultiGet(query []byte) []byte {
	keys := strings.Split(string(query), " ")
	if len(keys) == 0 {
		return []byte("ERR")
	}

	results := make([][]byte, len(keys))
	for i, key := range keys {
		if storage.KeyHasExpired(key) {
			storage.DeleteByKey(key)
			results[i] = []byte("Key not found")
			continue
		}

		data, err := storage.GetByKey(key)
		if err != nil {
			results[i] = []byte("Key not found")
			continue
		}

		results[i] = data
	}

	return parsing.JoinByteSlices(results, []byte("\n"))
}
