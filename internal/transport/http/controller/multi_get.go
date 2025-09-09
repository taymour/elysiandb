package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/taymour/elysiandb/internal/storage"
	"github.com/valyala/fasthttp"
)

type multiGetEntry struct {
	Key string  `json:"key"`
	Val *string `json:"value"`
}

func MultiGetController(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")

	keys := strings.Split(string(ctx.QueryArgs().Peek("keys")), ",")

	var results = make([]multiGetEntry, len(keys))
	for i, key := range keys {
		if storage.KeyHasExpired(key) {
			storage.DeleteByKey(key)

			results[i] = multiGetEntry{
				Key: key,
				Val: nil,
			}

			continue
		}

		data, err := storage.GetByKey(key)
		if err != nil {
			results[i] = multiGetEntry{
				Key: key,
				Val: nil,
			}

			continue
		}

		val := string(data)
		results[i] = multiGetEntry{
			Key: key,
			Val: &val,
		}
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = ctx.Write(jsonData)
}
