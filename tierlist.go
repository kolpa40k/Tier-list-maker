package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"tierlist-app/models"
	appstorage "tierlist-app/storage"
)

type TierListScreen struct {
	window        fyne.Window
	categories    *[]models.Category
	catIdx        int
	tlIdx         int
	tierList      models.TierList
	tierRows      []*TierRow
	topBar        *fyne.Container
	poolBox       *fyne.Container
	poolContainer *fyne.Container
	editMode      bool
	editBtn       *widget.Button
	content       fyne.CanvasObject
}

func NewTierListScreen(win fyne.Window, categories *[]models.Category, catIdx, tlIdx int, tl models.TierList) fyne.CanvasObject {
	ts := &TierListScreen{
		window: win, categories: categories,
		catIdx: catIdx, tlIdx: tlIdx, tierList: tl,
	}
	ts.buildUI()
	return ts.content
}

func (ts *TierListScreen) buildUI() {
	backBtn := widget.NewButton("← Назад", func() {
		ts.window.SetContent(
			NewCategoryScreen(ts.window, ts.categories, ts.catIdx, (*ts.categories)[ts.catIdx]),
		)
	})
	saveBtn := widget.NewButton("Сохранить", func() {
		ts.save()
		dialog.ShowInformation("Сохранено", "Тир-лист сохранён", ts.window)
	})
	ts.editBtn = widget.NewButton("Изменить", func() {
		ts.toggleEditMode()
	})
	ts.topBar = container.NewHBox(backBtn, saveBtn, ts.editBtn)

	tiersVBox := container.NewVBox()
	ts.tierRows = nil
	for i := range ts.tierList.Tiers {
		tier := &ts.tierList.Tiers[i]
		row := NewTierRow(tier, ts)
		ts.tierRows = append(ts.tierRows, row)
		tiersVBox.Add(row)
	}

	// Пул
	ts.poolBox = container.NewHBox()
	ts.rebuildPool()

	addImgBtn := widget.NewButton("+ Картинка", func() { ts.addImageToPool() })
	addImgBtn.Importance = widget.HighImportance
	poolHeader := widget.NewLabelWithStyle(
		"Не размещено:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true},
	)

	poolScroll := container.NewHScroll(ts.poolBox)
	poolScroll.SetMinSize(fyne.NewSize(800, 120))

	ts.poolContainer = container.NewBorder(
		container.NewVBox(poolHeader, addImgBtn),
		nil, nil, nil, poolScroll,
	)

	ts.content = container.NewBorder(
		ts.topBar, ts.poolContainer, nil, nil,
		container.NewVScroll(tiersVBox),
	)
}

func (ts *TierListScreen) toggleEditMode() {
	ts.editMode = !ts.editMode
	if ts.editMode {
		ts.editBtn.SetText("Готово")
	} else {
		ts.editBtn.SetText("Изменить")
	}
	for _, row := range ts.tierRows {
		row.Refresh()
	}
	ts.rebuildPool()
}

func (ts *TierListScreen) rebuildPool() {
	ts.poolBox.Objects = nil
	for i := range ts.tierList.Pool {
		item := ts.tierList.Pool[i]
		itemWidget := NewTierItemWidget(
			item,
			ts.editMode,
			func() { ts.deleteFromPool(item.ID) },
			func() { ts.renameInPool(item.ID) },
		)
		if ts.editMode {
			ts.poolBox.Add(itemWidget)
		} else {
			draggable := NewDraggableItem(itemWidget, item, nil, ts)
			ts.poolBox.Add(draggable)
		}
	}
	ts.poolBox.Refresh()
}

func (ts *TierListScreen) deleteFromPool(id string) {
	for i, it := range ts.tierList.Pool {
		if it.ID == id {
			ts.tierList.Pool = append(ts.tierList.Pool[:i], ts.tierList.Pool[i+1:]...)
			break
		}
	}
	ts.rebuildPool()
	ts.save()
}

func (ts *TierListScreen) renameInPool(id string) {
	var current *models.TierItem
	for i := range ts.tierList.Pool {
		if ts.tierList.Pool[i].ID == id {
			current = &ts.tierList.Pool[i]
			break
		}
	}
	if current == nil {
		return
	}

	entry := widget.NewEntry()
	entry.SetText(current.Label)
	entryBox := container.NewGridWrap(fyne.NewSize(360, 40), entry)

	dialog.ShowForm("Переименовать подпись", "Сохранить", "Отмена",
		[]*widget.FormItem{widget.NewFormItem("Подпись", entryBox)},
		func(ok bool) {
			if !ok || entry.Text == "" {
				return
			}
			current.Label = entry.Text
			ts.rebuildPool()
			ts.save()
		}, ts.window,
	)
}

func (ts *TierListScreen) addImageToPool() {
	dlg := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		path := reader.URI().Path()
		reader.Close()

		entry := widget.NewEntry()
		entry.SetPlaceHolder("Подпись (появится под картинкой)")
		entryBox := container.NewGridWrap(fyne.NewSize(360, 40), entry)

		dialog.ShowForm("Подпись картинки", "Добавить", "Отмена",
			[]*widget.FormItem{widget.NewFormItem("Подпись", entryBox)},
			func(ok bool) {
				if !ok {
					return
				}
				label := entry.Text
				if label == "" {
					label = path
				}
				item := models.TierItem{
					ID:        NewItemID(),
					Label:     label,
					ImagePath: path,
				}
				ts.tierList.Pool = append(ts.tierList.Pool, item)
				ts.rebuildPool()
				ts.save()
			}, ts.window,
		)
	}, ts.window)
	dlg.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	dlg.Show()
}

func (ts *TierListScreen) save() {
	(*ts.categories)[ts.catIdx].TierLists[ts.tlIdx] = ts.tierList
	if err := appstorage.SaveCategories(*ts.categories); err != nil {
		dialog.ShowError(err, ts.window)
	}
}

func (ts *TierListScreen) findDropTarget(absPos fyne.Position) (*TierRow, bool) {
	topBarH := ts.topBar.Size().Height
	for _, row := range ts.tierRows {
		rowY := topBarH + row.Position().Y
		if absPos.Y >= rowY && absPos.Y < rowY+row.Size().Height {
			return row, false
		}
	}
	poolY := ts.poolContainer.Position().Y
	if absPos.Y >= poolY {
		return nil, true
	}
	return nil, false
}

func (ts *TierListScreen) handleDrop(d *DraggableItem, pos fyne.Position) {
	if ts.editMode {
		return
	}
	target, isPool := ts.findDropTarget(pos)

	var sourceList *[]models.TierItem
	if d.sourceTier != nil {
		sourceList = &d.sourceTier.Items
	} else {
		sourceList = &ts.tierList.Pool
	}
	found := false
	for i, it := range *sourceList {
		if it.ID == d.item.ID {
			*sourceList = append((*sourceList)[:i], (*sourceList)[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return
	}

	if isPool {
		ts.tierList.Pool = append(ts.tierList.Pool, d.item)
	} else if target != nil {
		target.tier.Items = append(target.tier.Items, d.item)
	} else {
		*sourceList = append(*sourceList, d.item)
		return
	}

	for _, row := range ts.tierRows {
		row.Refresh()
	}
	ts.rebuildPool()
	ts.save()
}
