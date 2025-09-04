package boot

import (
	"time"

	"github.com/taymour/elysiandb/internal/storage"
)

func BootSaver() {
	go saveDBPeriodically()
}

func saveDBPeriodically() {
	for {
		storage.WriteToDB()
		time.Sleep(5 * time.Second)
	}
}
