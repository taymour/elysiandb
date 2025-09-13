package api

import (
	"encoding/json"

	"github.com/google/uuid"
	api_storage "github.com/taymour/elysiandb/internal/api"
	"github.com/valyala/fasthttp"
)

func CreateController(ctx *fasthttp.RequestCtx) {
	entity := ctx.UserValue("entity").(string)
	uuid := uuid.New().String()

	var data map[string]interface{}
	if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Invalid JSON")
		return
	}

	data["id"] = uuid

	api_storage.WriteEntity(entity, data)

	response, err := json.Marshal(data)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Error processing request")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(response)
}
