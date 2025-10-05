package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadConfig(path string) (GameConfig, error) {
	var config GameConfig

	var fileBytes, err = os.ReadFile(filepath.Join(path, "ember.json"))
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(fileBytes, &config); err != nil {
		return config, err
	}

	return config, err
}
