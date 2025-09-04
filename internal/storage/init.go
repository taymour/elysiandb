package storage

import (
	"maps"
	"fmt"
	"os"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/log"
)

var store *Store

func LoadDB() {
	cfg := globals.GetConfig()

	createFolder(cfg.Folder)
	createFile(cfg.Folder)

	store = newStore()

	data, err := ReadFromDB()
	if err != nil {
		log.Fatal("Error loading database:", err)
	}

	bytesData := make(map[string][]byte, len(data))
	maps.Copy(bytesData, data)
	store.FromMap(bytesData)
	store.saved.Store(true)
}

func createFolder(folder string) {
	if err := os.MkdirAll(folder, 0755); err != nil {
		log.Fatal("Error creating data folder:", err)
	}
}

func createFile(folder string) {
	filePath := folder + "/" + DataFile
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal("Error creating data file:", err)
		}
		file.WriteString("{}")
		defer file.Close()
	}
}

func GetByKey(key string) ([]byte, error) {
	if val, ok := store.get(key); ok {
		return val, nil
	}
	return nil, fmt.Errorf("key not found: %s", key)
}

func PutKeyValue(key string, value []byte) error {
	store.put(key, value)
	return nil
}

func DeleteByKey(key string) {
	store.del(key)
}

func ResetStore() {
	store.reset()
	log.Info("Store has been reset")
}

func Saved() bool {
	return store.saved.Load()
}
