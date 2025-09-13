package api_storage

import (
	"encoding/json"

	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/storage"
)

func WriteEntity(entity string, data map[string]interface{}) {
	key := globals.ApiSingleEntityKey(entity, data["id"].(string))

	jsonData, _ := json.Marshal(data)
	storage.PutKeyValue(key, jsonData)
}

func ReadlAllEntities(entity string) []map[string]interface{} {
	prefix := globals.ApiEntitiesAllKey(entity)
	data := storage.GetByWildcardKey(prefix)

	var results []map[string]interface{}
	for _, item := range data {
		var obj map[string]interface{}
		if err := json.Unmarshal(item, &obj); err == nil {
			results = append(results, obj)
		}
	}

	return results
}

func ReadEntityById(entity string, id string) map[string]interface{} {
	key := globals.ApiSingleEntityKey(entity, id)
	data, _ := storage.GetByKey(key)

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err == nil {
		return result
	}

	return nil
}

func DeleteEntityById(entity string, id string) {
	key := globals.ApiSingleEntityKey(entity, id)
	storage.DeleteByKey(key)
}

func DeleteAllEntities(entity string) {
	prefix := globals.ApiEntitiesAllKey(entity)
	storage.DeleteByWildcardKey(prefix)
}

func UpdateEntityById(entity string, id string, updated map[string]interface{}) map[string]interface{} {
	existing := ReadEntityById(entity, id)
	if existing == nil {
		return nil
	}

	for k, v := range updated {
		existing[k] = v
	}

	WriteEntity(entity, existing)
	return existing
}
