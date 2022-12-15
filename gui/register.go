package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sum-project/pass-locker2/db/repository"
)

func registerPage(a *App) fyne.CanvasObject {
	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Podaj login...")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Podaj hasło...")

	passwordConfirmEntry := widget.NewPasswordEntry()
	passwordConfirmEntry.SetPlaceHolder("Podaj hasło ponownie...")

	registerForm := widget.NewForm(
		widget.NewFormItem("", loginEntry),
		widget.NewFormItem("", passwordEntry),
		widget.NewFormItem("", passwordConfirmEntry),
	)
	registerForm.SubmitText = "Zarejestruj się"
	registerForm.OnSubmit = func() {
		user, err := a.userRepository.InsertUser(repository.InsertUserParams{
			Login:    loginEntry.Text,
			Password: passwordEntry.Text,
		})
		if err != nil {
			a.errorLog.Fatalln(err)
		}
		a.user = &user
		gui(a)
		a.infoLog.Println(user)
	}

	return container.NewVBox(registerForm)
}
