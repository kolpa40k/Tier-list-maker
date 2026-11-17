package ui

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/wierdling/gofiledialog"

	"tierlist-app/models"
	appstorage "tierlist-app/storage"
)

type TierListScreen struct {
	window       fyne.Window
	categories   *[]models.Category
	catIdx       int
	tlIdx        int
	tierList     models.TierList
	tierRows     []*TierRow
	topBar       *fyne.Container
	poolBox      *fyne.Container
	poolPanel    fyne.CanvasObject
	poolWrapper  *fyne.Container
	poolVisible  bool
	editMode     bool
	editBtn      *widget.Button
	deleteAllBtn *widget.Button
	content      fyne.CanvasObject
}

func NewTierListScreen(win fyne.Window, categories *[]models.Category, catIdx, tlIdx int, tl models.TierList) fyne.CanvasObject {
	ts := &TierListScreen{
		window: win, categories: categories,
		catIdx: catIdx, tlIdx: tlIdx, tierList: tl,
	}
	ts.tierList.RepairIDs()
	ts.buildUI()
	ts.setupDropHandler()
	ts.setupCloseHandler()
	return ts.content
}

func (ts *TierListScreen) buildUI() {
	// --- Верхняя панель ---
	backBtn := widget.NewButton("← Назад", func() {
		ts.window.SetOnDropped(nil)
		ts.window.SetContent(
			NewCategoryScreen(ts.window, ts.categories, ts.catIdx, (*ts.categories)[ts.catIdx]),
		)
	})
	saveBtn := widget.NewButton("Сохранить", func() {
		ts.save()
		dialog.ShowInformation("Сохранено", "Тир-лист сохранён", ts.window)
	})
	ts.editBtn = widget.NewButton("Изменить", func() { ts.toggleEditMode() })
	ts.topBar = container.NewHBox(backBtn, saveBtn, ts.editBtn)

	titleText := canvas.NewText(ts.tierList.Name, theme.Color(theme.ColorNameForeground))
	titleText.TextSize = 28 // РАЗМЕР ШРИФТА ТУТ!!!
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignCenter
	titleContainer := container.NewPadded(titleText)

	topSection := container.NewVBox(ts.topBar, titleContainer)

	// --- Тиры (вертикально) ---
	tiersVBox := container.NewVBox()
	ts.tierRows = nil
	for i := range ts.tierList.Tiers {
		tier := &ts.tierList.Tiers[i]
		row := NewTierRow(tier, ts)
		ts.tierRows = append(ts.tierRows, row)
		tiersVBox.Add(row)
	}

	// --- Основной контент ---
	mainContent := container.NewBorder(
		topSection, nil, nil, nil,
		container.NewVScroll(tiersVBox),
	)
	// --- FAB (круглая кнопка) в нижнем правом углу ---
	fab := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		ts.togglePool()
	})
	fab.Importance = widget.HighImportance
	fabBox := container.NewGridWrap(fyne.NewSize(56, 56), fab)
	fabRow := container.NewHBox(layout.NewSpacer(), fabBox)
	fabOverlay := container.NewBorder(
		nil, container.NewPadded(fabRow), nil, nil, nil,
	)

	// --- Боковая панель (скрытая) ---
	ts.poolBox = container.NewGridWrap(fyne.NewSize(100, 100))
	ts.rebuildPool()

	addImgBtn := widget.NewButton("+ Картинка", func() { ts.addImageToPool() })
	addImgBtn.Importance = widget.HighImportance

	ts.deleteAllBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		ts.deleteAllPoolItems()
	})
	ts.deleteAllBtn.Importance = widget.DangerImportance
	ts.updateDeleteAllVisibility()

	closeBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		ts.togglePool()
	})
	closeBtn.Importance = widget.LowImportance

	rightButtons := container.NewHBox(ts.deleteAllBtn, closeBtn)
	header := container.NewBorder(nil, nil, nil, rightButtons, addImgBtn)

	poolScroll := container.NewVScroll(ts.poolBox)
	inner := container.NewBorder(container.NewPadded(header), nil, nil, nil, poolScroll)
	innerPadded := container.NewPadded(inner)

	bg := canvas.NewRectangle(theme.Color(theme.ColorNameMenuBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameSeparator)
	bg.StrokeWidth = 1

	ts.poolPanel = container.NewStack(bg, innerPadded)

	ts.poolWrapper = container.NewWithoutLayout(ts.poolPanel)

	canvasSize := ts.window.Canvas().Size()
	if canvasSize.Width == 0 || canvasSize.Height == 0 {
		canvasSize = fyne.NewSize(900, 600)
	}
	ts.poolWrapper.Resize(canvasSize)
	panelW := ts.panelWidth(canvasSize)
	ts.poolPanel.Resize(fyne.NewSize(panelW, canvasSize.Height))
	ts.poolPanel.Move(fyne.NewPos(canvasSize.Width, 0))

	ts.content = container.NewStack(mainContent, fabOverlay, ts.poolWrapper)
}

func (ts *TierListScreen) panelWidth(size fyne.Size) float32 {
	w := size.Width / 4
	if w < 240 {
		w = 240
	}
	if w > 400 {
		w = 400
	}
	return w
}

func (ts *TierListScreen) togglePool() {
	canvasSize := ts.window.Canvas().Size()
	if canvasSize.Width == 0 || canvasSize.Height == 0 {
		canvasSize = fyne.NewSize(900, 600)
	}

	panelW := ts.panelWidth(canvasSize)
	ts.poolWrapper.Resize(canvasSize)
	ts.poolPanel.Resize(fyne.NewSize(panelW, canvasSize.Height))

	var startX, endX float32
	if ts.poolVisible {
		startX = canvasSize.Width - panelW
		endX = canvasSize.Width
	} else {
		startX = canvasSize.Width
		endX = canvasSize.Width - panelW
	}

	ts.poolPanel.Move(fyne.NewPos(startX, 0))
	ts.poolPanel.Refresh()

	anim := fyne.NewAnimation(280*time.Millisecond, func(d float32) {
		x := startX + (endX-startX)*d
		ts.poolPanel.Move(fyne.NewPos(x, 0))
		ts.poolPanel.Refresh()
	})
	anim.Start()

	ts.poolVisible = !ts.poolVisible
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

	if len(ts.tierList.Pool) == 0 {
		ts.poolBox.Layout = layout.NewVBoxLayout()
		ts.poolBox.Add(NewDropZone())
	} else {
		ts.poolBox.Layout = layout.NewGridWrapLayout(fyne.NewSize(100, 100))
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
	}

	ts.poolBox.Refresh()
	ts.updateDeleteAllVisibility()
}

func (ts *TierListScreen) deleteFromPool(id string) {
	for i, it := range ts.tierList.Pool {
		if it.ID == id {
			appstorage.DeleteImage(it.ImagePath)
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
			if !ok {
				return
			}
			current.Label = entry.Text
			ts.rebuildPool()
			ts.save()
		}, ts.window,
	)
}

func (ts *TierListScreen) addImageToPool() {
	err := gofiledialog.ShowOpen(func(paths []string, err error) {
		if err != nil {
			dialog.ShowError(err, ts.window)
			return
		}
		if len(paths) == 0 {
			return
		}

		for _, path := range paths {
			id := NewItemID()

			storedPath, err := appstorage.ImportImage(path, id)
			if err != nil {
				dialog.ShowError(
					fmt.Errorf("не удалось импортировать %s: %w", filepath.Base(path), err),
					ts.window,
				)
				continue
			}

			label := filepath.Base(path)
			item := models.TierItem{
				ID:        id,
				Label:     label,
				ImagePath: storedPath,
			}
			ts.tierList.Pool = append(ts.tierList.Pool, item)
		}

		ts.rebuildPool()
		ts.save()
	}, ts.window,
		gofiledialog.WithTitle("Выберите картинки"),
		gofiledialog.WithFilters(
			gofiledialog.Filter{Name: "Изображения", Extensions: []string{".png", ".jpg", ".jpeg", ".webp"}},
			gofiledialog.Filter{Name: "Все файлы"},
		),
		gofiledialog.WithMultiSelect(true),
	)
	if err != nil {
		dialog.ShowError(err, ts.window)
	}
}

func (ts *TierListScreen) save() {
	(*ts.categories)[ts.catIdx].TierLists[ts.tlIdx] = ts.tierList
	if err := appstorage.SaveCategories(*ts.categories); err != nil {
		dialog.ShowError(err, ts.window)
	}
}

func (ts *TierListScreen) findDropTarget(absPos fyne.Position) (*TierRow, bool) {
	if ts.poolVisible {
		p := ts.poolPanel.Position()
		s := ts.poolPanel.Size()
		if absPos.X >= p.X && absPos.X <= p.X+s.Width &&
			absPos.Y >= p.Y && absPos.Y <= p.Y+s.Height {
			return nil, true
		}
	}
	topBarH := ts.topBar.Size().Height
	for _, row := range ts.tierRows {
		rowY := topBarH + row.Position().Y
		if absPos.Y >= rowY && absPos.Y < rowY+row.Size().Height {
			return row, false
		}
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

func (ts *TierListScreen) updateDeleteAllVisibility() {
	if ts.deleteAllBtn == nil {
		return
	}
	if len(ts.tierList.Pool) == 0 {
		ts.deleteAllBtn.Hide()
	} else {
		ts.deleteAllBtn.Show()
	}
}

func (ts *TierListScreen) deleteAllPoolItems() {
	if len(ts.tierList.Pool) == 0 {
		return
	}
	count := len(ts.tierList.Pool)
	dialog.ShowConfirm(
		"Удалить все непривязанные картинки?",
		fmt.Sprintf(
			"Будет удалено %d картинок из панели.\n"+
				"Картинки, уже размещённые в тирах, останутся.",
			count,
		),
		func(ok bool) {
			if !ok {
				return
			}
			for _, it := range ts.tierList.Pool {
				appstorage.DeleteImage(it.ImagePath)
			}
			ts.tierList.Pool = nil
			ts.rebuildPool()
			ts.save()
		}, ts.window,
	)
}

func (ts *TierListScreen) setupDropHandler() {
	ts.window.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if len(uris) == 0 {
			return
		}

		var imageURIs []fyne.URI
		for _, u := range uris {
			ext := strings.ToLower(filepath.Ext(u.Path()))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp" {
				imageURIs = append(imageURIs, u)
			}
		}
		if len(imageURIs) == 0 {
			return
		}

		fyne.Do(func() {
			for _, u := range imageURIs {
				srcPath := u.Path()

				if runtime.GOOS == "windows" {
					srcPath = strings.TrimPrefix(srcPath, "/")
					srcPath = filepath.FromSlash(srcPath)
				}

				id := NewItemID()
				storedPath, err := appstorage.ImportImage(srcPath, id)
				if err != nil {
					log.Printf("не удалось импортировать %s: %v", srcPath, err)
					continue
				}

				item := models.TierItem{
					ID:        id,
					Label:     filepath.Base(srcPath),
					ImagePath: storedPath,
				}
				ts.tierList.Pool = append(ts.tierList.Pool, item)
			}

			ts.rebuildPool()
			ts.save()
		})
	})
}

func (ts *TierListScreen) setupCloseHandler() {
	ts.window.SetCloseIntercept(func() {
		fyne.CurrentApp().Quit()
	})
}
