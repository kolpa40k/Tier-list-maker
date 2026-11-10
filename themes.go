package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ThemePreset описывает одну цветовую схему.
type ThemePreset struct {
	Name   string
	Dark   bool // если true — базой служит тёмная тема Fyne
	Colors map[fyne.ThemeColorName]color.Color
}

// CustomTheme оборачивает базовую тему Fyne, подменяя часть цветов.
type CustomTheme struct {
	Base   fyne.Theme
	Colors map[fyne.ThemeColorName]color.Color
}

func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if c, ok := t.Colors[name]; ok {
		return c
	}
	return t.Base.Color(name, variant)
}
func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource    { return t.Base.Font(style) }
func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return t.Base.Icon(name) }
func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32       { return t.Base.Size(name) }

// Presets — список всех доступных тем.
// Ключевые цвета: Background (фон окна), Foreground (текст),
// Primary (акцент — кнопки HighImportance, фокус), Button (фон обычных кнопок),
// InputBackground (фон полей ввода), MenuBackground (фон диалогов и колонок тиров).
var Presets = []ThemePreset{
	{
		Name: "Классическая",
		Colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:      color.NRGBA{R: 245, G: 245, B: 245, A: 255},
			theme.ColorNameForeground:      color.NRGBA{R: 33, G: 33, B: 33, A: 255},
			theme.ColorNamePrimary:         color.NRGBA{R: 33, G: 150, B: 243, A: 255},
			theme.ColorNameButton:          color.NRGBA{R: 224, G: 224, B: 224, A: 255},
			theme.ColorNameInputBackground: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			theme.ColorNameMenuBackground:  color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		},
	},
	{
		Name: "Розовая",
		Colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:      color.NRGBA{R: 255, G: 228, B: 236, A: 255}, // #FFE4EC
			theme.ColorNameForeground:      color.NRGBA{R: 74, G: 26, B: 44, A: 255},    // #4A1A2C
			theme.ColorNamePrimary:         color.NRGBA{R: 233, G: 30, B: 99, A: 255},   // #E91E63
			theme.ColorNameButton:          color.NRGBA{R: 248, G: 187, B: 208, A: 255}, // #F8BBD0
			theme.ColorNameInputBackground: color.NRGBA{R: 255, G: 209, B: 220, A: 255}, // #FFD1DC
			theme.ColorNameMenuBackground:  color.NRGBA{R: 255, G: 240, B: 245, A: 255},
			theme.ColorNameSeparator:       color.NRGBA{R: 244, G: 143, B: 177, A: 255},
			theme.ColorNameHover:           color.NRGBA{R: 252, G: 200, B: 216, A: 255},
			theme.ColorNameSelection:       color.NRGBA{R: 248, G: 187, B: 208, A: 255},
		},
	},
	{
		Name: "Жёлтая",
		Colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:      color.NRGBA{R: 255, G: 248, B: 220, A: 255}, // #FFF8DC
			theme.ColorNameForeground:      color.NRGBA{R: 62, G: 39, B: 35, A: 255},    // #3E2723
			theme.ColorNamePrimary:         color.NRGBA{R: 249, G: 168, B: 37, A: 255},  // #F9A825
			theme.ColorNameButton:          color.NRGBA{R: 255, G: 224, B: 130, A: 255}, // #FFE082
			theme.ColorNameInputBackground: color.NRGBA{R: 255, G: 236, B: 179, A: 255}, // #FFECB3
			theme.ColorNameMenuBackground:  color.NRGBA{R: 255, G: 251, B: 230, A: 255},
			theme.ColorNameSeparator:       color.NRGBA{R: 255, G: 213, B: 79, A: 255},
			theme.ColorNameHover:           color.NRGBA{R: 255, G: 236, B: 179, A: 255},
			theme.ColorNameSelection:       color.NRGBA{R: 255, G: 224, B: 130, A: 255},
		},
	},
	{
		Name: "Тёмная",
		Dark: true,
		Colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:      color.NRGBA{R: 30, G: 30, B: 30, A: 255},    // #1E1E1E
			theme.ColorNameForeground:      color.NRGBA{R: 224, G: 224, B: 224, A: 255}, // #E0E0E0
			theme.ColorNamePrimary:         color.NRGBA{R: 187, G: 134, B: 252, A: 255}, // #BB86FC
			theme.ColorNameButton:          color.NRGBA{R: 45, G: 45, B: 45, A: 255},    // #2D2D2D
			theme.ColorNameInputBackground: color.NRGBA{R: 45, G: 45, B: 45, A: 255},    // #2D2D2D
			theme.ColorNameMenuBackground:  color.NRGBA{R: 40, G: 40, B: 40, A: 255},
			theme.ColorNameSeparator:       color.NRGBA{R: 60, G: 60, B: 60, A: 255},
			theme.ColorNameHover:           color.NRGBA{R: 55, G: 55, B: 55, A: 255},
			theme.ColorNameSelection:       color.NRGBA{R: 80, G: 60, B: 130, A: 255},
		},
	},
	{
		Name: "Синяя",
		Colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:      color.NRGBA{R: 227, G: 242, B: 253, A: 255}, // #E3F2FD
			theme.ColorNameForeground:      color.NRGBA{R: 13, G: 71, B: 161, A: 255},   // #0D47A1
			theme.ColorNamePrimary:         color.NRGBA{R: 25, G: 118, B: 210, A: 255},  // #1976D2
			theme.ColorNameButton:          color.NRGBA{R: 187, G: 222, B: 251, A: 255}, // #BBDEFB
			theme.ColorNameInputBackground: color.NRGBA{R: 187, G: 222, B: 251, A: 255}, // #BBDEFB
			theme.ColorNameMenuBackground:  color.NRGBA{R: 236, G: 247, B: 255, A: 255},
			theme.ColorNameSeparator:       color.NRGBA{R: 144, G: 202, B: 249, A: 255},
			theme.ColorNameHover:           color.NRGBA{R: 214, G: 236, B: 253, A: 255},
			theme.ColorNameSelection:       color.NRGBA{R: 187, G: 222, B: 251, A: 255},
		},
	},
	{
		Name: "Зелёная",
		Colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:      color.NRGBA{R: 232, G: 245, B: 233, A: 255}, // #E8F5E9
			theme.ColorNameForeground:      color.NRGBA{R: 27, G: 94, B: 32, A: 255},    // #1B5E20
			theme.ColorNamePrimary:         color.NRGBA{R: 67, G: 160, B: 71, A: 255},   // #43A047
			theme.ColorNameButton:          color.NRGBA{R: 200, G: 230, B: 201, A: 255}, // #C8E6C9
			theme.ColorNameInputBackground: color.NRGBA{R: 200, G: 230, B: 201, A: 255}, // #C8E6C9
			theme.ColorNameMenuBackground:  color.NRGBA{R: 241, G: 250, B: 242, A: 255},
			theme.ColorNameSeparator:       color.NRGBA{R: 165, G: 214, B: 167, A: 255},
			theme.ColorNameHover:           color.NRGBA{R: 214, G: 239, B: 216, A: 255},
			theme.ColorNameSelection:       color.NRGBA{R: 200, G: 230, B: 201, A: 255},
		},
	},
	{
		Name: "Фиолетовая",
		Colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:      color.NRGBA{R: 243, G: 229, B: 245, A: 255}, // #F3E5F5
			theme.ColorNameForeground:      color.NRGBA{R: 74, G: 20, B: 140, A: 255},   // #4A148C
			theme.ColorNamePrimary:         color.NRGBA{R: 123, G: 31, B: 162, A: 255},  // #7B1FA2
			theme.ColorNameButton:          color.NRGBA{R: 225, G: 190, B: 231, A: 255}, // #E1BEE7
			theme.ColorNameInputBackground: color.NRGBA{R: 225, G: 190, B: 231, A: 255}, // #E1BEE7
			theme.ColorNameMenuBackground:  color.NRGBA{R: 249, G: 240, B: 250, A: 255},
			theme.ColorNameSeparator:       color.NRGBA{R: 206, G: 147, B: 216, A: 255},
			theme.ColorNameHover:           color.NRGBA{R: 235, G: 210, B: 239, A: 255},
			theme.ColorNameSelection:       color.NRGBA{R: 225, G: 190, B: 231, A: 255},
		},
	},
}

func FindPreset(name string) ThemePreset {
	for _, p := range Presets {
		if p.Name == name {
			return p
		}
	}
	return Presets[0]
}

func ApplyPreset(app fyne.App, name string) string {
	preset := FindPreset(name)

	base := fyne.Theme(theme.LightTheme())
	if preset.Dark {
		base = theme.DarkTheme()
	}
	app.Settings().SetTheme(&CustomTheme{
		Base:   base,
		Colors: preset.Colors,
	})
	return preset.Name
}
