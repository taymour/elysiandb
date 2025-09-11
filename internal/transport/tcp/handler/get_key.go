package handler

import (
	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/stat"
	"github.com/taymour/elysiandb/internal/storage"
)

func HandleGet(query []byte) []byte {
	cfg := globals.GetConfig()
	if cfg.Stats.Enabled {
		stat.Stats.IncrementTotalRequests()
	}

	key := string(query)

	if storage.KeyHasExpired(key) {
		storage.DeleteByKey(key)

		if cfg.Stats.Enabled {
			stat.Stats.IncrementMisses()
		}

		return []byte("Key not found")
	}

	data, err := storage.GetByKey(key)
	if err != nil {
		if cfg.Stats.Enabled {
			stat.Stats.IncrementMisses()
		}

		return []byte("Key not found")
	}

	if cfg.Stats.Enabled {
		stat.Stats.IncrementHits()
	}

	return data
}
