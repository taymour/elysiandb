package controller

import (
	"net/http"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/stat"
	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

func DeleteKeyController(ctx *fasthttp.RequestCtx) {
	if globals.GetConfig().Stats.Enabled {
		stat.Stats.IncrementTotalRequests()
	}

	key := ctx.UserValue("key").(string)
	storage.DeleteByKey(key)
	ctx.SetStatusCode(http.StatusNoContent)
}
