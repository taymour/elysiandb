package storage

import (
	"encoding/json"
	"os"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/log"
)

func WriteToDB() {

	if Saved() {
		return
	}

	log.Info("Saving data to database...")

	cfg := globals.GetConfig()
	path := cfg.Folder + "/" + DataFile

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		log.Error("Error opening file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(store.ToMap()); err != nil {
		log.Error("Error encoding JSON:", err)
	}

	store.saved.Store(true)

	log.Info("Data saved to database.")
}
