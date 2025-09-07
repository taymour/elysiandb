package handler

import "github.com/taymour/elysiandb/internal/storage"

func HandleDelete(query []byte) []byte {
	storage.DeleteByKey(string(query))
	return []byte("OK")
}
