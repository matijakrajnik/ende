package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct{}

func (customTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if c == theme.ColorNameDisabled {
		return &color.RGBA{R: 175, G: 175, B: 175, A: 175}
	}
	return theme.DefaultTheme().Color(c, v)
}

func (customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (customTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (customTheme) Size(s fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(s)
}

func newCustomTheme() fyne.Theme {
	return &customTheme{}
}
