package helpers

import (
	"encoding/json"

	"github.com/spf13/cast"
)

func StructToStruct(s, k interface{}) error {

	body, err := json.Marshal(k)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &s)
	if err != nil {
		return err
	}

	return nil
}

func StructToMap(key string, data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, objectInterface := range cast.ToSlice(data) {
		var (
			objectMap = cast.ToStringMap(objectInterface)
			resultKey = cast.ToString(objectMap[key])
		)
		result[resultKey] = objectInterface
	}

	return result
}

func RemoveDuplicatesStrings(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}
