package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"tierlist-app/models"
	"tierlist-app/storage"
)

type CategoryScreen struct {
	window     fyne.Window
	categories *[]models.Category
	catIdx     int
	listBox    *fyne.Container
	editMode   bool
	editBtn    *widget.Button
	content    fyne.CanvasObject
}

func NewCategoryScreen(win fyne.Window, categories *[]models.Category, catIdx int, cat models.Category) fyne.CanvasObject {
	cs := &CategoryScreen{
		window:     win,
		categories: categories,
		catIdx:     catIdx,
		listBox:    container.NewVBox(),
	}

	cs.editBtn = widget.NewButton("Изменить", func() {
		cs.toggleEditMode()
	})

	createBtn := widget.NewButtonWithIcon("Создать тир-лист", nil, func() {
		cs.showCreateTierListDialog()
	})
	createBtn.Importance = widget.HighImportance

	backBtn := widget.NewButton("← Назад", func() {
		win.SetContent(NewHomeScreen(win, *categories))
	})

	top := container.NewHBox(backBtn, createBtn, cs.editBtn)
	cs.content = container.NewBorder(
		top, nil, nil, nil,
		container.NewVScroll(cs.listBox),
	)
	cs.rebuildList()
	return cs.content
}

func (cs *CategoryScreen) toggleEditMode() {
	cs.editMode = !cs.editMode
	if cs.editMode {
		cs.editBtn.SetText("Готово")
	} else {
		cs.editBtn.SetText("Изменить")
	}
	cs.rebuildList()
}

func (cs *CategoryScreen) rebuildList() {
	cs.listBox.Objects = nil
	cat := (*cs.categories)[cs.catIdx]
	for i := range cat.TierLists {
		i := i
		row := NewEditableRow(
			cat.TierLists[i].Name,
			cs.editMode,
			func() { cs.openTierList(i) },
			func() { cs.deleteTierList(i) },
			func() { cs.renameTierList(i) },
		)
		cs.listBox.Add(row)
	}
	cs.listBox.Refresh()
}

func (cs *CategoryScreen) renameTierList(idx int) {
	cat := &(*cs.categories)[cs.catIdx]
	entry := widget.NewEntry()
	entry.SetText(cat.TierLists[idx].Name)
	entryBox := container.NewGridWrap(fyne.NewSize(360, 40), entry)

	dialog.ShowForm("Переименовать тир-лист", "Сохранить", "Отмена",
		[]*widget.FormItem{widget.NewFormItem("Имя", entryBox)},
		func(ok bool) {
			if !ok || entry.Text == "" {
				return
			}
			cat.TierLists[idx].Name = entry.Text
			cs.rebuildList()
			if err := storage.SaveCategories(*cs.categories); err != nil {
				dialog.ShowError(err, cs.window)
			}
		}, cs.window,
	)
}

func (cs *CategoryScreen) openTierList(idx int) {
	cat := (*cs.categories)[cs.catIdx]
	tl := cat.TierLists[idx]
	cs.window.SetContent(
		NewTierListScreen(cs.window, cs.categories, cs.catIdx, idx, tl),
	)
}

func (cs *CategoryScreen) deleteTierList(idx int) {
	name := (*cs.categories)[cs.catIdx].TierLists[idx].Name
	dialog.ShowConfirm(
		"Удалить тир-лист?",
		"Тир-лист \""+name+"\" будет удалён.",
		func(ok bool) {
			if !ok {
				return
			}
			cat := &(*cs.categories)[cs.catIdx]
			cat.TierLists = append(cat.TierLists[:idx], cat.TierLists[idx+1:]...)
			cs.rebuildList()
			if err := storage.SaveCategories(*cs.categories); err != nil {
				dialog.ShowError(err, cs.window)
			}
		}, cs.window,
	)
}

func (cs *CategoryScreen) showCreateTierListDialog() {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Название тир-листа")

	entryBox := container.NewGridWrap(fyne.NewSize(360, 40), entry)

	dialog.ShowForm("Новый тир-лист", "Создать", "Отмена",
		[]*widget.FormItem{widget.NewFormItem("Имя", entryBox)},
		func(ok bool) {
			if !ok || entry.Text == "" {
				return
			}
			cat := &(*cs.categories)[cs.catIdx]
			cat.TierLists = append(cat.TierLists, models.NewTierList(entry.Text))
			cs.rebuildList()
			if err := storage.SaveCategories(*cs.categories); err != nil {
				dialog.ShowError(err, cs.window)
			}
		}, cs.window,
	)
}
