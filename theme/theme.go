package theme

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

//go:embed Lato.ttf
var font []byte

var myFont = &fyne.StaticResource{
	StaticName:    "Lato.ttf",
	StaticContent: font,
}

//go:embed logo.ico
var logo []byte

var MyLogo = &fyne.StaticResource{
	StaticName:    "logo.ico",
	StaticContent: logo,
}

type MyTheme struct {
	Theme string
}

var _ fyne.Theme = (*MyTheme)(nil)

func (m *MyTheme) Font(_ fyne.TextStyle) fyne.Resource {
	return myFont
}

func (m *MyTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameScrollBar {
		return 0
	}
	if n == theme.SizeNameScrollBarSmall {
		return 0
	}
	return theme.DefaultTheme().Size(n)
}

func (m *MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch m.Theme {
	case "Dark":
		v = theme.VariantDark
	case "Light":
		v = theme.VariantLight
	}
	return theme.DefaultTheme().Color(n, v)
}

func (m *MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
