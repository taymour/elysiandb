package handler

import (
	"strings"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/stat"
	"github.com/taymour/elysiandb/internal/storage"
	"github.com/taymour/elysiandb/internal/transport/tcp/parsing"
)

func HandleMultiGet(query []byte) []byte {
	cfg := globals.GetConfig()
	if cfg.Stats.Enabled {
		stat.Stats.IncrementTotalRequests()
	}

	keys := strings.Split(string(query), " ")
	if len(keys) == 0 {
		return []byte("ERR")
	}

	results := make([][]byte, len(keys))
	for i, key := range keys {
		if storage.KeyHasExpired(key) {
			storage.DeleteByKey(key)

			if cfg.Stats.Enabled {
				stat.Stats.IncrementMisses()
			}

			results[i] = []byte("Key not found")
			continue
		}

		data, err := storage.GetByKey(key)
		if err != nil {
			if cfg.Stats.Enabled {
				stat.Stats.IncrementMisses()
			}

			results[i] = []byte("Key not found")
			continue
		}

		results[i] = data

		if cfg.Stats.Enabled {
			stat.Stats.IncrementHits()
		}
	}

	return parsing.JoinByteSlices(results, []byte("\n"))
}
