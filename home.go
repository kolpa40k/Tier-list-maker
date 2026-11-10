package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"tierlist-app/models"
	"tierlist-app/storage"
)

type HomeScreen struct {
	window       fyne.Window
	categories   []models.Category
	listBox      *fyne.Container
	editMode     bool
	editBtn      *widget.Button
	currentTheme string
	content      fyne.CanvasObject
}

func NewHomeScreen(win fyne.Window, categories []models.Category) fyne.CanvasObject {
	h := &HomeScreen{
		window:     win,
		categories: categories,
		listBox:    container.NewVBox(),
	}

	if s, err := storage.LoadSettings(); err == nil && s != nil {
		h.currentTheme = s.ThemeName
	} else {
		h.currentTheme = "Классическая"
	}

	h.editBtn = widget.NewButton("Изменить", func() {
		h.toggleEditMode()
	})

	createBtn := widget.NewButtonWithIcon("Создать категорию", theme.ContentAddIcon(), func() {
		h.showCreateCategoryDialog()
	})
	createBtn.Importance = widget.HighImportance

	themeBtn := widget.NewButtonWithIcon("Тема", theme.ColorPaletteIcon(), func() {
		h.showThemeDialog()
	})

	top := container.NewHBox(createBtn, h.editBtn, themeBtn)
	h.content = container.NewBorder(
		top, nil, nil, nil,
		container.NewVScroll(h.listBox),
	)
	h.rebuildList()
	return h.content
}

func (h *HomeScreen) showThemeDialog() {
	names := make([]string, len(Presets))
	for i, p := range Presets {
		names[i] = p.Name
	}
	radio := widget.NewRadioGroup(names, nil)
	radio.Horizontal = false
	radio.SetSelected(h.currentTheme)

	radio.OnChanged = func(s string) {
		if s == "" {
			return
		}
		h.currentTheme = ApplyPreset(fyne.CurrentApp(), s)
		if err := storage.SaveSettings(&storage.Settings{ThemeName: h.currentTheme}); err != nil {
			dialog.ShowError(err, h.window)
		}
		h.rebuildList()
	}

	dialog.ShowCustom("Тема оформления", "Закрыть", radio, h.window)
}

func (h *HomeScreen) toggleEditMode() {
	h.editMode = !h.editMode
	if h.editMode {
		h.editBtn.SetText("Готово")
	} else {
		h.editBtn.SetText("Изменить")
	}
	h.rebuildList()
}

func (h *HomeScreen) rebuildList() {
	h.listBox.Objects = nil
	for i := range h.categories {
		i := i
		row := NewEditableRow(
			h.categories[i].Name,
			h.editMode,
			func() { h.openCategory(i) },
			func() { h.deleteCategory(i) },
			func() { h.renameCategory(i) },
		)
		h.listBox.Add(row)
	}
	h.listBox.Refresh()
}

func (h *HomeScreen) renameCategory(idx int) {
	entry := widget.NewEntry()
	entry.SetText(h.categories[idx].Name)
	entryBox := container.NewGridWrap(fyne.NewSize(360, 40), entry)

	dialog.ShowForm("Переименовать категорию", "Сохранить", "Отмена",
		[]*widget.FormItem{widget.NewFormItem("Имя", entryBox)},
		func(ok bool) {
			if !ok || entry.Text == "" {
				return
			}
			h.categories[idx].Name = entry.Text
			h.rebuildList()
			if err := storage.SaveCategories(h.categories); err != nil {
				dialog.ShowError(err, h.window)
			}
		}, h.window,
	)
}

func (h *HomeScreen) openCategory(idx int) {
	h.window.SetContent(
		NewCategoryScreen(h.window, &h.categories, idx, h.categories[idx]),
	)
}

func (h *HomeScreen) deleteCategory(idx int) {
	name := h.categories[idx].Name
	dialog.ShowConfirm(
		"Удалить категорию?",
		"Категория \""+name+"\" и все её тир-листы будут удалены.",
		func(ok bool) {
			if !ok {
				return
			}
			h.categories = append(h.categories[:idx], h.categories[idx+1:]...)
			h.rebuildList()
			if err := storage.SaveCategories(h.categories); err != nil {
				dialog.ShowError(err, h.window)
			}
		}, h.window,
	)
}

func (h *HomeScreen) showCreateCategoryDialog() {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Название категории")

	entryBox := container.NewGridWrap(fyne.NewSize(360, 40), entry)

	dialog.ShowForm("Новая категория", "Создать", "Отмена",
		[]*widget.FormItem{widget.NewFormItem("Имя", entryBox)},
		func(ok bool) {
			if !ok || entry.Text == "" {
				return
			}
			h.categories = append(h.categories, models.Category{Name: entry.Text})
			h.rebuildList()
			if err := storage.SaveCategories(h.categories); err != nil {
				dialog.ShowError(err, h.window)
			}
		}, h.window,
	)
}
