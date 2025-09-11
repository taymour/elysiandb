package controller

import (
	"fmt"
	"net/http"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/stat"
	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

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

		ctx.Error(fmt.Errorf("key '%s' not found", key).Error(), http.StatusNotFound)
		return
	}

	data, err := storage.GetByKey(key)
	if err != nil {
		if cfg.Stats.Enabled {
			stat.Stats.IncrementMisses()
		}

		ctx.Error(fmt.Errorf("key '%s' not found", key).Error(), http.StatusNotFound)
		return
	}

	if cfg.Stats.Enabled {
		stat.Stats.IncrementHits()
	}

	_, _ = ctx.Write(data)
}
