package controller

import (
	"fmt"
	"net/http"

	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

func GetKeyController(ctx *fasthttp.RequestCtx) {
	// Optionally set a default content type; if you store arbitrary bytes,
	// consider omitting or setting based on metadata.
	ctx.SetContentType("application/octet-stream")

	key := ctx.UserValue("key").(string)

	data, err := storage.GetByKey(key)
	if err != nil {
		ctx.Error(fmt.Errorf("key '%s' not found", key).Error(), http.StatusNotFound)
		return
	}

	// Write the raw bytes
	_, _ = ctx.Write(data)
}
