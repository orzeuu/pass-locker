package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/orzeuu/pass-locker/db/repository"
	"github.com/orzeuu/pass-locker/password"
)

func (a *App) addPasswordPage() fyne.CanvasObject {
	itemEntry := widget.NewEntry()
	itemEntry.SetPlaceHolder("Podaj item...")

	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Podaj login...")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Podaj hasło...")

	passwordConfirmEntry := widget.NewPasswordEntry()
	passwordConfirmEntry.SetPlaceHolder("Dodaj hasło...")

	addPasswordForm := widget.NewForm(
		widget.NewFormItem("", widget.NewLabel("Item")),
		widget.NewFormItem("", itemEntry),
		widget.NewFormItem("", widget.NewLabel("Login")),
		widget.NewFormItem("", loginEntry),
		widget.NewFormItem("", widget.NewLabel("Hasło")),
		widget.NewFormItem("", passwordEntry),
		widget.NewFormItem("", widget.NewLabel("Potwierdz hasło")),
		widget.NewFormItem("", passwordConfirmEntry),
	)

	addPasswordForm.OnSubmit = func() {
		pwdText, err := password.Encrypt([]byte(passwordEntry.Text), []byte(a.userPassword))
		if err != nil {
			a.errorLog.Fatalln(err)
		}
		pwd, err := a.passwordRepository.AddPassword(repository.AddPasswordParams{
			Item:     itemEntry.Text,
			Login:    loginEntry.Text,
			Password: string(pwdText),
			UserId:   a.user.ID,
		})
		if err != nil {
			a.errorLog.Fatalln(err)
		}
		a.infoLog.Println(pwd)
		a.gui()
	}
	addPasswordForm.SubmitText = "Dodaj hasło"

	return container.NewVBox(
		addPasswordForm,
	)
}
