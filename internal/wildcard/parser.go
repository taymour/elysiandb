package wildcard

func KeyContainsWildcard(key string) bool {
	for _, char := range key {
		if char == '*' || char == '?' {
			return true
		}
	}

	return false
}
