package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type homeTab struct {
	content *container.TabItem
}

func newHomeTab() *homeTab {
	h := &homeTab{}
	h.generateUI()
	return h
}

func (h *homeTab) generateUI() {
	style := fyne.TextStyle{Bold: true}

	h.content = container.NewTabItemWithIcon(
		"HOME",
		theme.HomeIcon(),
		container.NewCenter(
			container.New(
				layout.NewVBoxLayout(),
				widget.NewCard(
					"ENcode/DEcode", "", container.New(
						layout.NewVBoxLayout(),
						widget.NewLabelWithStyle("Base64 Encoder/Decoder", fyne.TextAlignCenter, style),
						widget.NewLabelWithStyle("URL Encoder/Decoder", fyne.TextAlignCenter, style),
					),
				),
				widget.NewSeparator(),
				widget.NewCard(
					"ENcrypt/DEcrypt", "", container.New(
						layout.NewVBoxLayout(),
						widget.NewLabelWithStyle("AES Encryption/Decryption", fyne.TextAlignCenter, style),
					),
				),
				widget.NewSeparator(),
				widget.NewCard(
					"+ Some Extras", "", container.New(
						layout.NewVBoxLayout(),
						widget.NewLabelWithStyle("MD5 Hashing Tool", fyne.TextAlignCenter, style),
						widget.NewLabelWithStyle("RNG (Random Number Generator)", fyne.TextAlignCenter, style),
					),
				),
			),
		),
	)
}
