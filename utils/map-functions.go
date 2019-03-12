package utils

func Filter(mapToFilter map[string]interface{}, f func(string, interface{}) bool) map[string]interface{} {
	filteredMap := make(map[string]interface{})

	for key, value := range mapToFilter {
		if f(key, value) {
			filteredMap[key] = value
		}
	}

	return filteredMap
}

// Create Keys

// Create Values
