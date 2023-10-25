package assets

import (
	"encoding/json"
	"os"
)

func LoadJSON(path string) (map[string]string, error) {
	val := make(map[string]string)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &val)
	if err != nil {
		return nil, err
	}

	return val, nil
}
