package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sum-project/pass-locker2/db/repository"
)

func (a *App) loginPage() fyne.CanvasObject {
	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Podaj login...")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Podaj hasło...")

	loginForm := widget.NewForm(
		widget.NewFormItem("", loginEntry),
		widget.NewFormItem("", passwordEntry),
	)
	loginForm.SubmitText = "Zaloguj się"
	loginForm.OnSubmit = func() {
		user, err := a.userRepository.GetUser(repository.GetUserParams{
			Login:    loginEntry.Text,
			Password: passwordEntry.Text,
		})
		if err != nil {
			a.errorLog.Fatalln(err)
		}
		a.user = &user
		a.userPassword = passwordEntry.Text
		a.gui()
		a.infoLog.Println(user)
	}

	return container.NewVBox(loginForm)
}
