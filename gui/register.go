package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func registerPage() fyne.CanvasObject {
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

	return container.NewVBox(registerForm)
}
