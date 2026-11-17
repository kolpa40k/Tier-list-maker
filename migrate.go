package storage

import (
	"os"
	"path/filepath"
	"strings"

	"tierlist-app/models"
)

func imagesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tierlist-app", "images"), nil
}

func isInside(path, dir string) bool {
	absPath, err1 := filepath.Abs(path)
	absDir, err2 := filepath.Abs(dir)
	if err1 != nil || err2 != nil {
		return false
	}
	rel, err := filepath.Rel(absDir, absPath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..") && rel != "."
}

func MigrateAllImages(categories []models.Category) (bool, error) {
	dir, err := imagesDir()
	if err != nil {
		return false, err
	}

	changed := false

	migrateItem := func(item *models.TierItem) {
		if item.ImagePath == "" {
			return
		}
		if item.ID == "" {
			return
		}
		if isInside(item.ImagePath, dir) {
			return
		}
		if _, err := os.Stat(item.ImagePath); err != nil {
			return
		}
		newPath, err := ImportImage(item.ImagePath, item.ID)
		if err != nil {
			return
		}
		item.ImagePath = newPath
		changed = true
	}

	for ci := range categories {
		for ti := range categories[ci].TierLists {
			tl := &categories[ci].TierLists[ti]
			for tierIdx := range tl.Tiers {
				for i := range tl.Tiers[tierIdx].Items {
					migrateItem(&tl.Tiers[tierIdx].Items[i])
				}
			}
			for i := range tl.Pool {
				migrateItem(&tl.Pool[i])
			}
		}
	}

	return changed, nil
}
