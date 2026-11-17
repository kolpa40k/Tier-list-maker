package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DropZone struct {
	widget.BaseWidget
}

func NewDropZone() *DropZone {
	d := &DropZone{}
	d.ExtendBaseWidget(d)
	return d
}

func (d *DropZone) CreateRenderer() fyne.WidgetRenderer {
	return &dropZoneRenderer{zone: d}
}

type dropZoneRenderer struct {
	zone  *DropZone
	text  *widget.Label
	lines []*canvas.Line
}

func (r *dropZoneRenderer) Layout(size fyne.Size) {
	if r.text == nil {
		r.text = widget.NewLabel("Перетащите или добавьте изображения")
		r.text.Alignment = fyne.TextAlignCenter
		r.text.Wrapping = fyne.TextWrapWord
		r.text.TextStyle = fyne.TextStyle{Bold: true}
	}

	padding := float32(20)

	r.text.Resize(fyne.NewSize(size.Width-padding*2, size.Height))
	actualHeight := r.text.MinSize().Height

	r.text.Resize(fyne.NewSize(size.Width-padding*2, actualHeight))
	r.text.Move(fyne.NewPos(padding, (size.Height-actualHeight)/2))

	// --- отрисовка пунктирной рамки ---
	r.lines = nil

	col := theme.Color(theme.ColorNameDisabled)
	dash, gap := float32(6), float32(6)
	step := dash + gap
	margin := float32(10) // отступы от краев ячейки

	drawLine := func(x1, y1, x2, y2 float32) {
		l := canvas.NewLine(col)
		l.StrokeWidth = 2
		l.Position1 = fyne.NewPos(x1, y1)
		l.Position2 = fyne.NewPos(x2, y2)
		r.lines = append(r.lines, l)
	}

	w := size.Width - margin*2
	h := size.Height - margin*2

	// верхняя грань
	for x := float32(0); x < w; x += step {
		length := dash
		if x+length > w {
			length = w - x
		}
		drawLine(margin+x, margin, margin+x+length, margin)
	}
	// нижняя грань
	for x := float32(0); x < w; x += step {
		length := dash
		if x+length > w {
			length = w - x
		}
		drawLine(margin+x, margin+h, margin+x+length, margin+h)
	}
	// левая грань
	for y := float32(0); y < h; y += step {
		length := dash
		if y+length > h {
			length = h - y
		}
		drawLine(margin, margin+y, margin, margin+y+length)
	}
	// правая грань
	for y := float32(0); y < h; y += step {
		length := dash
		if y+length > h {
			length = h - y
		}
		drawLine(margin+w, margin+y, margin+w, margin+y+length)
	}
}

func (r *dropZoneRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (r *dropZoneRenderer) Refresh() {
	if r.text != nil {
		r.text.Refresh()
	}
	for _, l := range r.lines {
		l.Refresh()
	}
}

func (r *dropZoneRenderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{r.text}
	for _, l := range r.lines {
		objs = append(objs, l)
	}
	return objs
}

func (r *dropZoneRenderer) Destroy() {}
