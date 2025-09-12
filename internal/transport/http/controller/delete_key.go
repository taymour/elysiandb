package controller

import (
	"net/http"
	"net/url"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/stat"
	"github.com/taymour/elysiandb/internal/storage"
	"github.com/taymour/elysiandb/internal/wildcard"
	"github.com/valyala/fasthttp"
)

func DeleteKeyController(ctx *fasthttp.RequestCtx) {
	if globals.GetConfig().Stats.Enabled {
		stat.Stats.IncrementTotalRequests()
	}

	key := ctx.UserValue("key").(string)

	if dec, err := url.PathUnescape(key); err == nil {
		key = dec
	}

	if wildcard.KeyContainsWildcard(key) {
		storage.DeleteByWildcardKey(key)
	} else {
		storage.DeleteByKey(key)
	}

	ctx.SetStatusCode(http.StatusNoContent)
}
