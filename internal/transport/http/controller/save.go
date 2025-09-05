package controller

import (
	"net/http"

	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

func SaveController(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusNoContent)

	go storage.WriteToDB()
}
