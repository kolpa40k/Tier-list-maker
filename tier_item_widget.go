package ui

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/draw"

	"tierlist-app/models"
)

func loadImageResource(path string) fyne.Resource {
	if path == "" {
		return nil
	}
	tryPaths := []string{path}
	if len(path) > 0 && (path[0] == '/' || path[0] == '\\') {
		tryPaths = append(tryPaths, path[1:])
	}
	tryPaths = append(tryPaths, filepath.FromSlash(path))

	for _, p := range tryPaths {
		if data, err := os.ReadFile(p); err == nil && len(data) > 0 {
			return fyne.NewStaticResource(filepath.Base(p), data)
		}
	}
	fmt.Printf("[TierItemWidget] FAILED to open image, tried: %v\n", tryPaths)
	return nil
}

type TierItemWidget struct {
	widget.BaseWidget
	item     models.TierItem
	image    *canvas.Image
	editMode bool
	onDelete func()
	onRename func()
}

func NewTierItemWidget(item models.TierItem, editMode bool, onDelete, onRename func()) *TierItemWidget {
	w := &TierItemWidget{
		item:     item,
		editMode: editMode,
		onDelete: onDelete,
		onRename: onRename,
	}

	img := createThumbnail(item.ImagePath, 200, 200)
	if img == nil {
		img = canvas.NewImageFromResource(theme.BrokenImageIcon())
	}
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(100, 100))
	w.image = img

	w.ExtendBaseWidget(w)
	w.Resize(fyne.NewSize(100, 100))
	return w
}

func (w *TierItemWidget) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (w *TierItemWidget) CreateRenderer() fyne.WidgetRenderer {
	var content fyne.CanvasObject = w.image

	// 👇 Полоску и текст создаем ТОЛЬКО если подпись не пустая
	if w.item.Label != "" {
		strip := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 170})

		display := w.item.Label
		r := []rune(display)
		if len(r) > 14 {
			display = string(r[:14]) + "…"
		}
		labelText := canvas.NewText(display, color.White)
		labelText.TextSize = 12
		labelText.TextStyle = fyne.TextStyle{Bold: true}

		labelStack := container.NewStack(strip, container.NewCenter(labelText))
		labelBar := container.NewGridWrap(fyne.NewSize(100, 22), labelStack)
		bottomOverlay := container.NewBorder(nil, labelBar, nil, nil, nil)

		content = container.NewStack(w.image, bottomOverlay)
	}

	if w.editMode {
		renameBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			if w.onRename != nil {
				w.onRename()
			}
		})
		renameBtn.Importance = widget.LowImportance

		delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			if w.onDelete != nil {
				w.onDelete()
			}
		})
		delBtn.Importance = widget.DangerImportance

		buttons := container.NewHBox(renameBtn, delBtn)
		overlay := container.NewBorder(buttons, nil, nil, nil, nil)

		border := canvas.NewRectangle(color.Transparent)
		border.StrokeColor = color.RGBA{R: 220, A: 255}
		border.StrokeWidth = 2

		content = container.NewStack(content, overlay, border)
	}

	return widget.NewSimpleRenderer(content)
}

func createThumbnail(path string, width, height int) *canvas.Image {
	file, err := os.Open(path)
	if err != nil {
		fyne.LogError("Failed to open image file", err)
		return nil
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		fyne.LogError("Failed to decode image", err)
		return nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	return canvas.NewImageFromImage(dst)
}
