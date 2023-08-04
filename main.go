package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const (
	AppTitle = "ENDE"
)

var (
	myApp      fyne.App
	mainWindow fyne.Window
)

func main() {
	myApp = app.New()
	myApp.Settings().SetTheme(newCustomTheme())
	mainWindow = myApp.NewWindow(AppTitle)
	mainWindow.SetContent(generateMainUI())
	mainWindow.Resize(fyne.Size{Width: 800, Height: 500})
	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()
}

func generateMainUI() fyne.CanvasObject {
	return container.NewMax(container.NewAppTabs(
		newHomeTab().content,
		newBase64Tab().content,
		newUrlTab().content,
		newAesTab().content,
		newMd5Tab().content,
		newRngTab().content,
	))
}
