package storage

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"tierlist-app/models"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
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

func ImportImage(srcPath string, id string) (string, error) {
	dir, err := imagesDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("открыть исходник: %w", err)
	}
	defer srcFile.Close()

	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return "", fmt.Errorf("декодировать изображение: %w", err)
	}

	const targetSize = 200
	srcBounds := srcImg.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	scaleW := float64(targetSize) / float64(srcW)
	scaleH := float64(targetSize) / float64(srcH)
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}

	newW := int(float64(srcW) * scale)
	newH := int(float64(srcH) * scale)

	offsetX := (targetSize - newW) / 2
	offsetY := (targetSize - newH) / 2

	dstImg := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))

	draw.Draw(dstImg, dstImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	dstRect := image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH)
	draw.ApproxBiLinear.Scale(dstImg, dstRect, srcImg, srcBounds, draw.Over, nil)

	dstPath := filepath.Join(dir, id+".png")

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("создать файл: %w", err)
	}
	defer dstFile.Close()

	if err := png.Encode(dstFile, dstImg); err != nil {
		return "", fmt.Errorf("сохранить миниатюру: %w", err)
	}

	return dstPath, nil
}

func DeleteImage(path string) {
	if path == "" {
		return
	}
	dir, err := imagesDir()
	if err != nil {
		return
	}

	if !isInside(path, dir) {
		return
	}
	_ = os.Remove(path)
}
