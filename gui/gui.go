package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type App struct {
	Win fyne.Window
	App fyne.App
}

func Start() {
	a := initApp()

	a.Win.Resize(fyne.NewSize(a.Win.Canvas().Size().Width, a.Win.Canvas().Size().Height))
	a.Win.CenterOnScreen()
	a.Win.SetMaster()
	a.Win.ShowAndRun()
}

func initApp() *App {
	a := app.NewWithID("pass-locker")
	w := a.NewWindow("Pass Locker")

	return &App{
		Win: w,
		App: a,
	}
}
