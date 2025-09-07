package boot

import (
	"time"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/storage"
)

func BootSaver() {
	go saveDBPeriodically(
		time.Duration(
			globals.GetConfig().Store.FlushIntervalSeconds,
		) * time.Second,
	)
}

func saveDBPeriodically(interval time.Duration) {
	for {
		storage.WriteToDB()
		time.Sleep(interval)
	}
}
