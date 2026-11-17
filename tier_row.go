package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"tierlist-app/models"
	appstorage "tierlist-app/storage"
)

type TierRow struct {
	widget.BaseWidget
	tier     *models.Tier
	itemsBox *fyne.Container
	nameBtn  *TierNameButton
	screen   *TierListScreen
}

func NewTierRow(tier *models.Tier, screen *TierListScreen) *TierRow {
	tr := &TierRow{tier: tier, screen: screen}
	tr.itemsBox = container.NewHBox()
	tr.nameBtn = NewTierNameButton(tier.Name, screen.editMode, func() {
		tr.editName()
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
	tr.nameBtn.SetData(tr.tier.Name, tr.screen.editMode)
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
			appstorage.DeleteImage(it.ImagePath)
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
			if !ok {
				return
			}
			current.Label = entry.Text
			tr.Refresh()
			tr.screen.save()
		}, tr.screen.window,
	)
}

// editName — диалог переименования тира с валидацией длины (≤30) и счётчиком.
func (tr *TierRow) editName() {
	entry := widget.NewEntry()
	entry.SetText(tr.tier.Name)
	entryBox := container.NewGridWrap(fyne.NewSize(360, 40), entry)

	hint := widget.NewLabel("")
	hint.Alignment = fyne.TextAlignTrailing

	var dlg *dialog.CustomDialog

	saveBtn := widget.NewButton("Сохранить", func() {
		name := strings.TrimSpace(entry.Text)
		if name == "" {
			dialog.ShowError(fmt.Errorf("название не может быть пустым"), tr.screen.window)
			return
		}
		if utf8.RuneCountInString(name) > 30 {
			dialog.ShowError(
				fmt.Errorf("название не может быть длиннее 30 символов (сейчас %d)",
					utf8.RuneCountInString(name)),
				tr.screen.window,
			)
			return
		}
		tr.tier.Name = name
		tr.Refresh()
		tr.screen.save()
		dlg.Hide()
	})
	saveBtn.Importance = widget.HighImportance

	updateCounter := func(s string) {
		n := utf8.RuneCountInString(s)
		if n > 30 {
			hint.SetText(fmt.Sprintf("%d / 30 — превышен лимит", n))
			saveBtn.Disable()
		} else {
			hint.SetText(fmt.Sprintf("%d / 30", n))
			saveBtn.Enable()
		}
	}
	entry.OnChanged = updateCounter
	updateCounter(entry.Text)

	content := container.NewVBox(
		widget.NewLabel("Имя тира (максимум 30 символов):"),
		entryBox,
		hint,
		container.NewHBox(layout.NewSpacer(), saveBtn),
	)

	dlg = dialog.NewCustom("Переименовать тир", "Отмена", content, tr.screen.window)
	dlg.Resize(fyne.NewSize(420, 220))
	dlg.Show()
}
