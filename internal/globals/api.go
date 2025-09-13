package globals

import "fmt"

const (
	ApiEntityPattern       = "api:entity:%s"
	ApiEntitiesPattern     = "api:entity:%s:*"
	ApiSingleEntityPattern = "api:entity:%s:id:%s"
)

func ApiEntityKey(entity string) string {
	return fmt.Sprintf(ApiEntityPattern, entity)
}

func ApiEntitiesAllKey(entity string) string {
	return fmt.Sprintf(ApiEntitiesPattern, entity)
}

func ApiSingleEntityKey(entity string, id string) string {
	return fmt.Sprintf(ApiSingleEntityPattern, entity, id)
}
