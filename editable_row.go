package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type EditableRow struct {
	widget.BaseWidget
	text     string
	editMode bool
	onTap    func()
	onDelete func()
	onRename func()
}

func NewEditableRow(text string, editMode bool, onTap, onDelete, onRename func()) *EditableRow {
	r := &EditableRow{
		text: text, editMode: editMode,
		onTap: onTap, onDelete: onDelete, onRename: onRename,
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *EditableRow) CreateRenderer() fyne.WidgetRenderer {
	var textBtn *widget.Button
	if r.editMode {
		textBtn = widget.NewButton(r.text, nil)
		textBtn.Importance = widget.LowImportance
	} else {
		textBtn = widget.NewButton(r.text, r.onTap)
		textBtn.Importance = widget.HighImportance
	}
	textBtn.Alignment = widget.ButtonAlignLeading

	if !r.editMode {
		return widget.NewSimpleRenderer(textBtn)
	}

	renameBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), r.onRename)
	renameBtn.Importance = widget.LowImportance

	delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), r.onDelete)
	delBtn.Importance = widget.DangerImportance

	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.RGBA{R: 220, A: 255}
	border.StrokeWidth = 2

	buttons := container.NewHBox(renameBtn, delBtn)
	content := container.NewBorder(nil, nil, nil, buttons, textBtn)
	stack := container.NewStack(border, content)
	return widget.NewSimpleRenderer(stack)
}
