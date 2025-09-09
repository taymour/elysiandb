package routing

import (
	"github.com/fasthttp/router"
	"github.com/taymour/elysiandb/internal/transport/http/controller"
)

func RegisterRoutes(r *router.Router) {
	r.GET("/health", controller.HealthController)

	r.GET("/kv/mget", controller.MultiGetController)
	r.GET("/kv/{key}", controller.GetKeyController)
	r.PUT("/kv/{key}", controller.PutKeyController)
	r.DELETE("/kv/{key}", controller.DeleteKeyController)

	r.POST("/save", controller.SaveController)

	r.POST("/reset", controller.ResetController)
}
