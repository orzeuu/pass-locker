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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

type App struct {
	Win fyne.Window
	App fyne.App

	user *repository.User

	userRepository     repository.UserRepository
	passwordRepository repository.PasswordRepository

	errorLog *log.Logger
	infoLog  *log.Logger

	pwd            *widget.Entry
	pwdEntropy     binding.String
	pwdOptionsBind binding.BoolList
	lengthBind     binding.Float

	passwords     [][]interface{}
	passwordTable *widget.Table
}

func Start() {
	a := initApp()

	gui(a)

	a.Win.Resize(fyne.NewSize(a.Win.Canvas().Size().Width, a.Win.Canvas().Size().Height))
	a.Win.CenterOnScreen()
	a.Win.SetMaster()
	a.Win.ShowAndRun()
}

func gui(a *App) {
	tabs := container.NewAppTabs()

	if a.user != nil {
		tabs.Append(container.NewTabItemWithIcon("Generowanie hasła", theme.DocumentIcon(), container.NewPadded(generatorWindow(a))))
		tabs.Append(container.NewTabItemWithIcon("Lista haseł", theme.ListIcon(), container.NewPadded(passwordListPage(a))))
		tabs.Append(container.NewTabItemWithIcon("Dodaj hasło", theme.ContentAddIcon(), container.NewPadded(addPasswordPage(a))))
		tabs.Append(container.NewTabItemWithIcon("Ustawienia", theme.DocumentIcon(), container.NewPadded(settingsWindow())))
	}

	if a.user == nil {
		tabs.Append(container.NewTabItemWithIcon("Logowanie", theme.LoginIcon(), container.NewPadded(loginPage(a))))
		tabs.Append(container.NewTabItemWithIcon("Rejestracja", theme.HomeIcon(), container.NewPadded(registerPage(a))))
	}

	tabs.OnSelected = func(t *container.TabItem) {
		t.Content.Refresh()
	}

	a.Win.SetContent(tabs)
}

func initApp() *App {
	a := app.NewWithID("pass-locker")

	t := a.Preferences().StringWithFallback("Theme", "Light")
	a.Settings().SetTheme(&theme2.MyTheme{Theme: t})
	a.SetIcon(theme2.MyLogo)

	w := a.NewWindow("Pass Locker")

	db, err := initSQLite(a.Storage().RootURI().Path())
	if err != nil {
		panic(err)
	}

	return &App{
		Win: w,
		App: a,

		errorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),

		userRepository:     repository.NewUserRepository(db),
		passwordRepository: repository.NewPasswordRepository(db),
	}
}

func initSQLite(appPath string) (*gorm.DB, error) {
	path := appPath + "/sql.db"

	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&repository.User{})
	db.AutoMigrate(&repository.Password{})

	return db, nil
}
