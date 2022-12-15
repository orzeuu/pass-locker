package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sum-project/pass-locker2/password"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type Password struct {
	Numbers   bool
	Lowercase bool
	Uppercase bool
	Symbol    bool
	Length    uint
}

func generatorWindow(a *App) fyne.CanvasObject {
	number := a.App.Preferences().BoolWithFallback("Number", false)
	lowercase := a.App.Preferences().BoolWithFallback("Lowercase", false)
	uppercase := a.App.Preferences().BoolWithFallback("Uppercase", false)
	symbol := a.App.Preferences().BoolWithFallback("Symbol", false)

	length := a.App.Preferences().IntWithFallback("Length", 0)
	lengthBind := binding.NewFloat()
	_ = lengthBind.Set(float64(length))

	pwdOptionsBind := binding.NewBoolList()
	_ = pwdOptionsBind.Set([]bool{number, lowercase, uppercase, symbol})

	slide := widget.NewSliderWithData(0, 64, lengthBind)
	slide.Step = 1
	lengthText := widget.NewLabelWithData(binding.FloatToStringWithFormat(lengthBind, "Długość hasła: %0.0f"))

	buttons := container.NewGridWithColumns(
		4,
		widget.NewButton("8", func() {
			_ = lengthBind.Set(8)
		}),
		widget.NewButton("16", func() {
			_ = lengthBind.Set(16)
		}),
		widget.NewButton("32", func() {
			_ = lengthBind.Set(32)
		}),
		widget.NewButton("64", func() {
			_ = lengthBind.Set(64)
		}))

	lengthLabel := container.NewGridWithColumns(2, container.New(layout.NewFormLayout(), lengthText, slide), buttons)

	pwd := widget.NewEntry()
	a.pwd = pwd

	pwdEntropy := binding.NewString()
	pwdEntropyLabel := widget.NewLabelWithData(pwdEntropy)
	pwdEntropyText := canvas.NewText("Siła hasła ", nil)
	a.pwdEntropy = pwdEntropy

	slide.OnChanged = func(f float64) {
		_ = lengthBind.Set(f)
		a.App.Preferences().SetInt("Length", int(f))
	}

	lengthBind.AddListener(binding.NewDataListener(func() {
		pwdSetText(a)
	}))

	a.lengthBind = lengthBind
	a.pwdOptionsBind = pwdOptionsBind

	NumberCheck := widgetCheck(a, "Numer", "Number", number)
	LowercaseCheck := widgetCheck(a, "Małe litery", "Lowercase", lowercase)
	UppercaseCheck := widgetCheck(a, "Duże litery", "Uppercase", uppercase)
	SymbolCheck := widgetCheck(a, "Symbole", "Symbol", symbol)

	copyBtn := widget.NewButtonWithIcon("Kopiuj", theme.ContentCopyIcon(), func() {
		a.Win.Clipboard().SetContent(pwd.Text)
	})
	copyBtn.Importance = widget.HighImportance

	updateBth := widget.NewButtonWithIcon("Generuj", theme.ViewRefreshIcon(), func() {
		pwdSetText(a)
	})

	resetBth := widget.NewButtonWithIcon("Reset", theme.CancelIcon(), func() {
		pwd.SetText("")
		slide.SetValue(0)
		NumberCheck.SetChecked(false)
		LowercaseCheck.SetChecked(false)
		UppercaseCheck.SetChecked(false)
		SymbolCheck.SetChecked(false)
	})

	opButtons := container.New(layout.NewGridLayout(3), copyBtn, updateBth, resetBth)

	checklists := container.New(layout.NewGridLayout(4), NumberCheck, LowercaseCheck, UppercaseCheck, SymbolCheck)

	content := container.NewVBox(pwd, container.New(layout.NewFormLayout(), pwdEntropyText, pwdEntropyLabel), lengthLabel, checklists, opButtons)

	a.Win.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyRight, fyne.KeyDown:
			if slide.Value < slide.Max {
				_ = lengthBind.Set(slide.Value + slide.Step)
			}
		case fyne.KeyLeft, fyne.KeyUp:
			if slide.Value > slide.Min {
				_ = lengthBind.Set(slide.Value - slide.Step)
			}
		case fyne.KeyF5:
			pwdSetText(a)
		}
	})

	return content
}

func widgetCheck(a *App, label, key string, checked bool) *widget.Check {
	check := widget.NewCheck(label, func(b bool) {})
	check.SetChecked(checked)
	check.OnChanged = func(b bool) {
		a.App.Preferences().SetBool(key, b)
		number := a.App.Preferences().BoolWithFallback("Number", false)
		lowercase := a.App.Preferences().BoolWithFallback("Lowercase", false)
		uppercase := a.App.Preferences().BoolWithFallback("Uppercase", false)
		symbol := a.App.Preferences().BoolWithFallback("Symbol", false)
		_ = a.pwdOptionsBind.Set([]bool{number, lowercase, uppercase, symbol})

		pwdSetText(a)
	}

	return check
}

func genPwd(p *Password) string {
	config := password.Config{}
	config.IncludeNumbers = p.Numbers
	config.IncludeLowercaseLetters = p.Lowercase
	config.IncludeUppercaseLetters = p.Uppercase
	config.IncludeSymbols = p.Symbol

	if !p.Numbers && !p.Lowercase && !p.Uppercase && !p.Symbol {
		config.IncludeNumbers = true
	}

	config.Length = p.Length

	g := password.NewGenerator(&config)
	pwd, _ := g.Generate()

	return pwd
}

func pwdLevel(pwd string) string {
	e := passwordvalidator.GetEntropy(pwd)
	switch {
	case e < 20:
		return "bardzo słabe"
	case e < 40:
		return "słabe"
	case e < 60:
		return "średnie"
	case e < 80:
		return "mocne"
	case e >= 80:
		return "bardzo mocne"
	default:
		return ""
	}
}

func pwdSetText(a *App) {
	pwdOptions, _ := a.pwdOptionsBind.Get()
	length, _ := a.lengthBind.Get()

	if length > 0 {
		if !pwdOptions[0] && !pwdOptions[1] && !pwdOptions[2] && !pwdOptions[3] {
			a.pwd.SetText("")
			_ = a.pwdEntropy.Set("")
		} else {
			pwd := genPwd(&Password{pwdOptions[0], pwdOptions[1], pwdOptions[2], pwdOptions[3], uint(length)})
			a.pwd.SetText(pwd)
			_ = a.pwdEntropy.Set(pwdLevel(pwd))
		}
	} else {
		a.pwd.SetText("")
		_ = a.pwdEntropy.Set("")
	}
}
