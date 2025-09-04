package controller

import (
	"net/http"

	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

func PutKeyController(ctx *fasthttp.RequestCtx) {
	key := ctx.UserValue("key").(string)

	// fasthttp's body buffer is reused; ALWAYS copy before storing.
	body := ctx.PostBody()
	buf := make([]byte, len(body))
	copy(buf, body)

	if err := storage.PutKeyValue(key, buf); err != nil {
		ctx.Error("Failed to store key-value pair", http.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(http.StatusNoContent)
}
