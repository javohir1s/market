package helpers

import (
	"encoding/json"
	"os"
)

func Read(path string) ([]interface{}, error) {

	body, err := os.ReadFile(path)
	if err != nil {

		return nil, err
	}

	var response []interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func Write(path string, data []interface{}) error {

	body, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, body, os.ModePerm)
	if err != nil {
		return err
	}

	return err
}
