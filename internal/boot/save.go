package boot

import (
	"time"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/storage"
)

func BootSaver() {
	d := time.Duration(globals.GetConfig().Store.FlushIntervalSeconds) * time.Second
	if d <= 0 {
		return
	}

	go saveDBPeriodically(d)
}

func saveDBPeriodically(interval time.Duration) {
	for {
		storage.WriteToDB()
		time.Sleep(interval)
	}
}
