package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"tierlist-app/storage"
	"tierlist-app/ui"
)

func main() {
	myApp := app.New()

	settings, err := storage.LoadSettings()
	if err != nil || settings == nil {
		settings = &storage.Settings{ThemeName: "Классическая"}
	}
	ui.ApplyPreset(myApp, settings.ThemeName)

	myWindow := myApp.NewWindow("Tier List Maker")
	myWindow.Resize(fyne.NewSize(900, 600))

	categories, err := storage.LoadCategories()
	if err != nil {
		categories = nil
	}

	myWindow.SetContent(ui.NewHomeScreen(myWindow, categories))
	myWindow.ShowAndRun()
}
