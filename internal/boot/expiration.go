package boot

import (
	"time"

	"github.com/taymour/elysiandb/internal/storage"
)

func BootExpirationHandler() {
	go checkExpirationsPeriodically()
}

func checkExpirationsPeriodically() {
	for {
		storage.CleanExpiratedKeys(time.Now().Unix())
		time.Sleep(1 * time.Second)
	}
}
