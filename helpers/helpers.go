package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func CreateProject(path string, name string) error {
	var fullPath = filepath.Join(path, name)

	var _, statErr = os.Stat(fullPath)
	if !os.IsNotExist(statErr) {
		return errors.New("A folder with this name already exists in this location")
	}

	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return err
	}

	var _, err = os.Create(filepath.Join(fullPath, "ember.json"))
	if err != nil {
		return fmt.Errorf("%s\nError occured during the creation of the %v.ember file", err, name)
	}
	return nil
}

func IsValidProject(path string) error {
	var configPath = filepath.Join(path, "ember.json")
	var _, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return errors.New("Not a valid ember project!")
	}

	return nil
}
