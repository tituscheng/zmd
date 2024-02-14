package gui

import (
	"fmt"
	"zmd/cache"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const (
	WINDOWS_WIDTH  = 1000
	WINDOWS_HEIGHT = 600
)

var (
	mediaLocationString = binding.NewString()
)

func Run() {
	cache.Init()
	a := app.New()
	w := a.NewWindow("ZTM player")

	tabs := container.NewAppTabs(
		container.NewTabItem("Settings", SetupSettingsContainer()),
		container.NewTabItem("Media", widget.NewLabel("World!")),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(WINDOWS_WIDTH, WINDOWS_HEIGHT))
	w.ShowAndRun()
}

func SetupSettingsContainer() *fyne.Container {
	updateEntry := widget.NewEntry()
	updateEntry.PlaceHolder = "/Users/"

	updateMediaLocationLabel()

	mediaLocationLabel := widget.NewLabelWithData(mediaLocationString)
	mediaLocationRefreshButton := widget.NewButton("Refresh", func() {
		updateMediaLocationLabel()
	})
	mediaLocationContainer := container.NewBorder(nil, nil, nil, mediaLocationRefreshButton, mediaLocationLabel)

	updateButton := widget.NewButton("Update", func() {
		cache.WriteAndUpdateMediaLocation(updateEntry.Text)
		updateEntry.SetText("")
		updateMediaLocationLabel()
	})

	updateContainer := container.NewBorder(nil, nil, nil, updateButton, updateEntry)
	mainContainer := container.NewVBox(mediaLocationContainer, updateContainer)
	return mainContainer
}

func updateMediaLocationLabel() {
	fileCache := GetMediaLocation()
	mediaLocationString.Set(generateMediaLabel(fileCache.MediaLocation))
}

func generateMediaLabel(path string) string {
	if len(path) == 0 {
		return "Media folder: location not found"
	}
	return fmt.Sprintf("Media folder: %s", path)
}

func GetMediaLocation() cache.Cache {
	return cache.GetCache()
}
