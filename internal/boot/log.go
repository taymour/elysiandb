package boot

import (
	"time"

	"github.com/taymour/elysiandb/internal/log"
)

func BootLogger() {
	go WriteLogsPeriodically()
}

func WriteLogsPeriodically() {
	for {
		if len(log.Logs) > 0 {
			log.WriteLogs()
		}

		time.Sleep(5 * time.Second)
	}
}
