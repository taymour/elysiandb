package globals

import "github.com/taymour/elysiandb/internal/configuration"

var (
	cfg *configuration.Config
)

func SetConfig(c *configuration.Config) {
	cfg = c
}

func GetConfig() *configuration.Config {
	return cfg
}
