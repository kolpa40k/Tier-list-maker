package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"tierlist-app/models"
)

const dataFile = "categories.json"

func dataPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tierlist-app", dataFile), nil
}

func LoadCategories() ([]models.Category, error) {
	path, err := dataPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []models.Category{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var categories []models.Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func SaveCategories(categories []models.Category) error {
	path, err := dataPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
