package main

import (
	"fmt"

	"github.com/taymour/elysiandb/internal/boot"
	"github.com/taymour/elysiandb/internal/configuration"
	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/log"
)

func main() {
	fmt.Println(`
   ╔══════════════════════════════════════╗
   ║                                      ║
   ║      Welcome to ElysianDB            ║
   ║  A modern, lightweight KV datastore  ║
   ║                                      ║
   ╚══════════════════════════════════════╝
	`)

	cfg, err := configuration.LoadConfig("elysian.yaml")
	if err != nil {
		log.Error("Error loading config:", err)
		return
	}

	globals.SetConfig(cfg)

	log.Info("Using data folder: ", globals.GetConfig().Store.Folder)

	boot.InitDB()

	log.Info("Ready to serve your key-value needs with elegance.")

	if cfg.Server.HTTP.Enabled {
		go boot.StartHTTP()
	}

	if cfg.Server.TCP.Enabled {
		go boot.InitTCP()
	}

	// Block forever
	select {}
}
