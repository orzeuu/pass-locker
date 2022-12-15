package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func loginPage() fyne.CanvasObject {
	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Podaj login...")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Podaj hasło...")

	loginForm := widget.NewForm(
		widget.NewFormItem("", loginEntry),
		widget.NewFormItem("", passwordEntry),
	)
	loginForm.SubmitText = "Zaloguj się"

	return container.NewVBox(loginForm)
}
