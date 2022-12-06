package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	theme2 "github.com/sum-project/pass-locker2/theme"
)

type App struct {
	Win fyne.Window
	App fyne.App
}

func Start() {
	a := initApp()

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Generowanie hasła", theme.DocumentIcon(), container.NewPadded()),
		container.NewTabItemWithIcon("Ustawienia", theme.DocumentIcon(), container.NewPadded(settingsWindow())),
	)

	tabs.OnSelected = func(t *container.TabItem) {
		t.Content.Refresh()
	}

	a.Win.SetContent(tabs)
	a.Win.Resize(fyne.NewSize(a.Win.Canvas().Size().Width, a.Win.Canvas().Size().Height))
	a.Win.CenterOnScreen()
	a.Win.SetMaster()
	a.Win.ShowAndRun()
}

func initApp() *App {
	a := app.NewWithID("pass-locker")

	t := a.Preferences().StringWithFallback("Theme", "Light")
	a.Settings().SetTheme(&theme2.MyTheme{Theme: t})
	a.SetIcon(theme2.MyLogo)

	w := a.NewWindow("Pass Locker")

	return &App{
		Win: w,
		App: a,
	}
}
