package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sum-project/pass-locker2/db/repository"
	theme2 "github.com/sum-project/pass-locker2/theme"
)

type App struct {
	Win fyne.Window
	App fyne.App

	userRepository     repository.UserRepository
	passwordRepository repository.PasswordRepository

	pwd            *widget.Entry
	pwdEntropy     binding.String
	pwdOptionsBind binding.BoolList
	lengthBind     binding.Float
}

func Start() {
	a := initApp()

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Logowanie", theme.LoginIcon(), container.NewPadded(loginPage())),
		container.NewTabItemWithIcon("Rejestracja", theme.HomeIcon(), container.NewPadded(registerPage())),
		container.NewTabItemWithIcon("Generowanie hasła", theme.DocumentIcon(), container.NewPadded(generatorWindow(a))),
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
