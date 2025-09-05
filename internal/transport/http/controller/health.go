package controller

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

func HealthController(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}
