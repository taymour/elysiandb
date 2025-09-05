package boot

import (
	"github.com/taymour/elysiandb/internal/storage"
)

func InitDB() {
	storage.LoadDB()
	BootSaver()
	BootExpirationHandler()
	BootLogger()
}
