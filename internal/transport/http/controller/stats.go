package controller

import (
	"github.com/taymour/elysiandb/internal/stat"
	"github.com/valyala/fasthttp"
)

func StatsController(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	_, _ = ctx.Write([]byte(stat.Stats.ToJson()))
}
