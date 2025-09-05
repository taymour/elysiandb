package handler

import "github.com/taymour/elysiandb/internal/storage"

func HandleGet(query []byte) []byte {
	key := string(query)

	if storage.KeyHasExpired(key) {
		storage.DeleteByKey(key)
		return []byte("Key not found")
	}

	data, err := storage.GetByKey(key)
	if err != nil {
		return []byte("Key not found")
	}

	return data
}
