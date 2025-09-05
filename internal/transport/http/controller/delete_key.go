package controller

import (
	"net/http"

	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

func DeleteKeyController(ctx *fasthttp.RequestCtx) {
	key := ctx.UserValue("key").(string)
	storage.DeleteByKey(key)
	ctx.SetStatusCode(http.StatusNoContent)
}
