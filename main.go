package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"

	"tierlist-app/storage"
	"tierlist-app/ui"
)

func main() {
	myApp := app.New()
	settings, err := storage.LoadSettings()
	if err != nil {
		log.Printf("не удалось загрузить настройки: %v", err)
	}
	if settings == nil || settings.ThemeName == "" {
		settings = &storage.Settings{ThemeName: "Классическая"}
	}
	ui.ApplyPreset(myApp, settings.ThemeName)

	myWindow := myApp.NewWindow("Tier List Maker")
	myWindow.Resize(fyne.NewSize(900, 600))

	categories, loadErr := storage.LoadCategories()
	if loadErr != nil {
		log.Printf("не удалось загрузить категории: %v", loadErr)
	}

	myWindow.SetContent(ui.NewHomeScreen(myWindow, categories))

	if loadErr != nil {
		myApp.Lifecycle().SetOnStarted(func() {
			dialog.ShowError(
				fmt.Errorf(
					"Не удалось загрузить сохранённые категории.\n"+
						"Приложение запущено с пустым списком.\n\n%v",
					loadErr,
				),
				myWindow,
			)
		})
	}

	myWindow.ShowAndRun()
}
