package handler

import (
	"github.com/taymour/elysiandb/internal/storage"
	"github.com/taymour/elysiandb/internal/transport/tcp/parser"
)

func HandleSet(query []byte) []byte {
	k, v := parser.FirstWordBytes(query)

	key := string(k)
	val := make([]byte, len(v))
	copy(val, v)

	if err := storage.PutKeyValue(key, val); err != nil {
		return []byte("ERR")
	}
	return []byte("OK")
}
