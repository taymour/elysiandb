package controller

import (
	"encoding/json"
	"net/http"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/stat"
	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

type getEntry struct {
	Key string  `json:"key"`
	Val *string `json:"value"`
}

func GetKeyController(ctx *fasthttp.RequestCtx) {
	cfg := globals.GetConfig()
	if cfg.Stats.Enabled {
		stat.Stats.IncrementTotalRequests()
	}

	ctx.SetContentType("text/plain; charset=utf-8")

	key := ctx.UserValue("key").(string)

	if storage.KeyHasExpired(key) {
		storage.DeleteByKey(key)

		if cfg.Stats.Enabled {
			stat.Stats.IncrementMisses()
		}

		jsonData, _ := json.Marshal(getEntry{
			Key: key,
			Val: nil,
		})
		_, _ = ctx.Write(jsonData)
		ctx.SetStatusCode(http.StatusNotFound)

		return
	}

	data, err := storage.GetByKey(key)
	if err != nil {
		if cfg.Stats.Enabled {
			stat.Stats.IncrementMisses()
		}

		jsonData, _ := json.Marshal(getEntry{
			Key: key,
			Val: nil,
		})
		_, _ = ctx.Write(jsonData)
		ctx.SetStatusCode(http.StatusNotFound)

		return
	}

	if cfg.Stats.Enabled {
		stat.Stats.IncrementHits()
	}

	valStr := string(data)
	jsonData, _ := json.Marshal(getEntry{
		Key: key,
		Val: &valStr,
	})

	_, _ = ctx.Write(jsonData)
}
