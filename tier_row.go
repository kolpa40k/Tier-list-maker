package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"tierlist-app/models"
)

type TierRow struct {
	widget.BaseWidget
	tier     *models.Tier
	itemsBox *fyne.Container
	nameBtn  *widget.Button
	screen   *TierListScreen
}

func NewTierRow(tier *models.Tier, screen *TierListScreen) *TierRow {
	tr := &TierRow{tier: tier, screen: screen}
	tr.itemsBox = container.NewHBox()

	tr.nameBtn = widget.NewButton(tier.Name, func() {
		if tr.screen.editMode {
			tr.editName()
		}
	})

	tr.ExtendBaseWidget(tr)
	tr.Refresh()
	return tr
}

func (tr *TierRow) Refresh() {
	tr.itemsBox.Objects = nil
	for i := range tr.tier.Items {
		item := tr.tier.Items[i]
		itemWidget := NewTierItemWidget(
			item,
			tr.screen.editMode,
			func() { tr.deleteItem(item.ID) },
			func() { tr.renameItem(item.ID) },
		)
		if tr.screen.editMode {
			tr.itemsBox.Add(itemWidget)
		} else {
			draggable := NewDraggableItem(itemWidget, item, tr.tier, tr.screen)
			tr.itemsBox.Add(draggable)
		}
	}
	tr.itemsBox.Refresh()

	tr.nameBtn.SetText(tr.tier.Name)
	if tr.screen.editMode {
		tr.nameBtn.Importance = widget.WarningImportance
		tr.nameBtn.SetIcon(theme.DocumentCreateIcon())
	} else {
		tr.nameBtn.Importance = widget.HighImportance
		tr.nameBtn.SetIcon(nil)
	}
	tr.nameBtn.Refresh()
}

func (tr *TierRow) CreateRenderer() fyne.WidgetRenderer {
	nameBox := container.NewGridWrap(fyne.NewSize(120, 100), tr.nameBtn)
	itemsScroll := container.NewHScroll(tr.itemsBox)
	itemsScroll.SetMinSize(fyne.NewSize(600, 100))
	return widget.NewSimpleRenderer(
		container.NewBorder(nil, nil, nameBox, nil, itemsScroll),
	)
}

func (tr *TierRow) deleteItem(id string) {
	for i, it := range tr.tier.Items {
		if it.ID == id {
			tr.tier.Items = append(tr.tier.Items[:i], tr.tier.Items[i+1:]...)
			break
		}
	}
	tr.Refresh()
	tr.screen.save()
}

func (tr *TierRow) renameItem(id string) {
	var current *models.TierItem
	for i := range tr.tier.Items {
		if tr.tier.Items[i].ID == id {
			current = &tr.tier.Items[i]
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
			tr.Refresh()
			tr.screen.save()
		}, tr.screen.window,
	)
}

func (tr *TierRow) editName() {
	entry := widget.NewEntry()
	entry.SetText(tr.tier.Name)
	entryBox := container.NewGridWrap(fyne.NewSize(360, 40), entry)

	dialog.ShowForm("Переименовать тир", "Сохранить", "Отмена",
		[]*widget.FormItem{widget.NewFormItem("Имя", entryBox)},
		func(ok bool) {
			if !ok || entry.Text == "" {
				return
			}
			tr.tier.Name = entry.Text
			tr.Refresh()
			tr.screen.save()
		}, tr.screen.window,
	)
}
