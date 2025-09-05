package controller

import (
	"fmt"
	"net/http"

	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

func GetKeyController(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain; charset=utf-8")

	key := ctx.UserValue("key").(string)

	if storage.KeyHasExpired(key) {
		storage.DeleteByKey(key)
		ctx.Error(fmt.Errorf("key '%s' not found", key).Error(), http.StatusNotFound)
		return
	}

	data, err := storage.GetByKey(key)
	if err != nil {
		ctx.Error(fmt.Errorf("key '%s' not found", key).Error(), http.StatusNotFound)
		return
	}

	_, _ = ctx.Write(data)
}
