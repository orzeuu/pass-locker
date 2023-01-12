package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/orzeuu/pass-locker/db/repository"
	theme2 "github.com/orzeuu/pass-locker/theme"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

type App struct {
	Win fyne.Window
	App fyne.App

	user         *repository.User
	userPassword string

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

	a.gui()

	a.Win.Resize(fyne.NewSize(600, 600))
	a.Win.CenterOnScreen()
	a.Win.SetMaster()
	a.Win.ShowAndRun()
}

func (a *App) gui() {
	tabs := container.NewAppTabs()

	if a.user != nil {
		tabs.Append(container.NewTabItemWithIcon("Generowanie hasła", theme.DocumentIcon(), container.NewPadded(a.generatorWindow())))
		tabs.Append(container.NewTabItemWithIcon("Lista haseł", theme.ListIcon(), container.NewPadded(a.passwordListPage())))
		tabs.Append(container.NewTabItemWithIcon("Dodaj hasło", theme.ContentAddIcon(), container.NewPadded(a.addPasswordPage())))
		tabs.Append(container.NewTabItemWithIcon("Ustawienia", theme.DocumentIcon(), container.NewPadded(a.settingsWindow())))
	}

	if a.user == nil {
		tabs.Append(container.NewTabItemWithIcon("Logowanie", theme.LoginIcon(), container.NewPadded(a.loginPage())))
		tabs.Append(container.NewTabItemWithIcon("Rejestracja", theme.HomeIcon(), container.NewPadded(a.registerPage())))
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
	test := a.Storage().RootURI().Path()
	test = test + test
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
