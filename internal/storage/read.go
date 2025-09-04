package storage

import (
	"encoding/json"
	"io"
	"os"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/log"
)

func ReadFromDB() (map[string][]byte, error) {
	cfg := globals.GetConfig()

	file, err := os.Open(cfg.Folder + "/" + DataFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	log.Info("Successfully Opened "+DataFile, " for reading.")

	data := make(map[string][]byte)

	byteValue, _ := io.ReadAll(file)
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return nil, err
	}

	return data, nil
}
