package boot

import (
	"fmt"

	"github.com/fasthttp/router"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/log"
	"github.com/taymour/elysiandb/internal/routing"
	"github.com/valyala/fasthttp"
)

func StartHTTP() {
	cfg := globals.GetConfig()

	addr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)

	r := router.New()
	routing.RegisterRoutes(r)

	srv := &fasthttp.Server{
		Handler:               r.Handler,
		Name:                  "ElysianDB",
		DisableKeepalive:      false,
		TCPKeepalive:          true,
		ReadBufferSize:        4096,
		WriteBufferSize:       4096,
		ReduceMemoryUsage:     true,
		LogAllErrors:          false,
		NoDefaultServerHeader: true,
	}

	log.Info("ElysianDB HTTP listening on http://", addr)
	if err := srv.ListenAndServe(addr); err != nil {
		log.Fatal("server error: ", err)
	}

	log.WriteLogs()
}
