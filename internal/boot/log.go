package boot

import (
	"time"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/log"
)

func BootLogger() {
	go WriteLogsPeriodically(
		time.Duration(
			globals.GetConfig().Log.FlushIntervalSeconds,
		) * time.Second,
	)
}

func WriteLogsPeriodically(interval time.Duration) {
	for {
		log.WriteLogs()
		time.Sleep(interval)
	}
}
