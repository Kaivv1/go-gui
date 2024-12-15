package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct {
	fyne.Theme
}

func newCustomTheme() fyne.Theme {
	return &CustomTheme{Theme: theme.DefaultTheme()}
}

func (t *CustomTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, theme.VariantDark)
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return BASE_TEXT
	}
	return t.Theme.Size(name) + 4
}
