// ui/draggable_item.go
package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"tierlist-app/models"
)

func NewItemID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type DraggableItem struct {
	widget.BaseWidget
	inner      fyne.CanvasObject
	item       models.TierItem
	sourceTier *models.Tier
	screen     *TierListScreen
	lastPos    fyne.Position
	ghost      *widget.PopUp
}

func NewDraggableItem(inner fyne.CanvasObject, item models.TierItem, sourceTier *models.Tier, screen *TierListScreen) *DraggableItem {
	d := &DraggableItem{inner: inner, item: item, sourceTier: sourceTier, screen: screen}
	d.ExtendBaseWidget(d)
	return d
}

func (d *DraggableItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(d.inner)
}

func (d *DraggableItem) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (d *DraggableItem) Dragged(ev *fyne.DragEvent) {
	d.lastPos = ev.AbsolutePosition

	if d.ghost == nil {
		res := loadImageResource(d.item.ImagePath)
		ghostImg := canvas.NewImageFromResource(res)
		ghostImg.FillMode = canvas.ImageFillContain
		ghostImg.SetMinSize(fyne.NewSize(100, 100))
		ghostImg.Translucency = 0.5

		d.ghost = widget.NewPopUp(ghostImg, d.screen.window.Canvas())
		d.ghost.Show()
	}

	d.ghost.Move(fyne.NewPos(ev.AbsolutePosition.X-50, ev.AbsolutePosition.Y-50))
}

func (d *DraggableItem) DragEnd() {
	// Прячем и удаляем "призрак".
	if d.ghost != nil {
		d.ghost.Hide()
		d.ghost = nil
	}

	if d.screen != nil {
		d.screen.handleDrop(d, d.lastPos)
	}
}
