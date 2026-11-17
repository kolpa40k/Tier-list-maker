package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TierNameButton struct {
	widget.BaseWidget
	text     string
	editMode bool
	onTap    func()

	label *widget.Label
	bg    *canvas.Rectangle
}

func NewTierNameButton(text string, editMode bool, onTap func()) *TierNameButton {
	b := &TierNameButton{text: text, editMode: editMode, onTap: onTap}

	b.label = widget.NewLabel(text)
	b.label.Wrapping = fyne.TextWrapWord
	b.label.Alignment = fyne.TextAlignCenter
	b.label.TextStyle = fyne.TextStyle{Bold: true}

	b.bg = canvas.NewRectangle(b.bgColor())
	b.bg.CornerRadius = 6

	b.ExtendBaseWidget(b)
	return b
}

func (b *TierNameButton) bgColor() color.Color {
	if b.editMode {
		return theme.Color(theme.ColorNameWarning)
	}
	return theme.Color(theme.ColorNamePrimary)
}

func (b *TierNameButton) SetData(text string, editMode bool) {
	b.text = text
	b.editMode = editMode
	b.label.SetText(text)
	b.bg.FillColor = b.bgColor()
	b.bg.Refresh()
}

func (b *TierNameButton) CreateRenderer() fyne.WidgetRenderer {
	tap := widget.NewButton("", func() {
		if b.editMode && b.onTap != nil {
			b.onTap()
		}
	})
	tap.Importance = widget.LowImportance

	centered := container.NewVBox(
		layout.NewSpacer(),
		b.label,
		layout.NewSpacer(),
	)

	content := container.NewStack(
		b.bg,
		container.NewPadded(centered),
		tap,
	)
	return widget.NewSimpleRenderer(content)
}
