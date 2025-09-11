package controller

import (
	"net/http"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/stat"
	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

func PutKeyController(ctx *fasthttp.RequestCtx) {
	if globals.GetConfig().Stats.Enabled {
		stat.Stats.IncrementTotalRequests()
	}

	key := ctx.UserValue("key").(string)
	ttl := ctx.QueryArgs().GetUintOrZero("ttl")

	body := ctx.PostBody()
	buf := make([]byte, len(body))
	copy(buf, body)

	var err error
	if ttl > 0 {
		err = storage.PutKeyValueWithTTL(key, buf, ttl)
	} else {
		err = storage.PutKeyValue(key, buf)
	}

	if err != nil {
		ctx.Error("Failed to store key-value pair", http.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(http.StatusNoContent)
}
