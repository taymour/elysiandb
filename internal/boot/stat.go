package boot

import (
	"time"

	"github.com/taymour/elysiandb/internal/stat"
)

func BootStats() {
	stat.Init()
	go collect()
}

func collect() {
	for {
		stat.Stats.IncrementUptimeSeconds()
		time.Sleep(1 * time.Second)
	}
}
