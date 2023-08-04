package main

import (
	"encoding/base64"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type base64Tab struct {
	content *container.TabItem

	decEntry    *widget.Entry
	decEntryOut *widget.Entry
	decButton   *widget.Button

	encEntry    *widget.Entry
	encEntryOut *widget.Entry
	encButton   *widget.Button
}

func newBase64Tab() *base64Tab {
	b64 := &base64Tab{}
	b64.generateUI()
	b64.setButtonsHandlers()
	return b64
}

func (b64 *base64Tab) generateUI() {
	b64.content = container.NewTabItem(
		"BASE 64",
		container.New(
			layout.NewGridLayout(2),
			b64.generateDecodeUI(),
			b64.generateEncodeUI(),
		),
	)
}

func (b64 *base64Tab) generateDecodeUI() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Base64 decoder", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true})
	outLabel := widget.NewLabelWithStyle("Decoded string:", fyne.TextAlignLeading, fyne.TextStyle{})

	b64.decEntry = widget.NewMultiLineEntry()
	b64.decEntryOut = widget.NewMultiLineEntry()
	b64.decButton = widget.NewButton("DECODE", func() {})

	b64.decEntry.SetPlaceHolder("Enter Base64 encoded string here")
	b64.decEntryOut.Disable()

	outLabelCont := container.NewBorder(nil, outLabel, nil, nil)
	return container.NewBorder(title, nil, nil, nil, container.New(layout.NewGridLayout(1), b64.decEntry, b64.decButton, outLabelCont, b64.decEntryOut))
}

func (b64 *base64Tab) generateEncodeUI() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Base64 encoder", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true})
	outLabel := widget.NewLabelWithStyle("Encoded string:", fyne.TextAlignLeading, fyne.TextStyle{})

	b64.encEntry = widget.NewMultiLineEntry()
	b64.encEntryOut = widget.NewMultiLineEntry()
	b64.encButton = widget.NewButton("ENCODE", func() {})

	b64.encEntry.SetPlaceHolder("Enter string here to encode it in Base64")
	b64.encEntryOut.Disable()

	outLabelCont := container.NewBorder(nil, outLabel, nil, nil)
	return container.NewBorder(title, nil, nil, nil, container.New(layout.NewGridLayout(1), b64.encEntry, b64.encButton, outLabelCont, b64.encEntryOut))
}

func (b64 *base64Tab) setButtonsHandlers() {
	// Set DECODE button handler
	b64.decButton.OnTapped = func() {
		if b64.decEntry.Text == "" {
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(b64.decEntry.Text)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}
		b64.decEntryOut.SetText(string(decoded))
	}

	// Set ENCODE button handler
	b64.encButton.OnTapped = func() {
		if b64.encEntry.Text == "" {
			return
		}
		encoded := base64.StdEncoding.EncodeToString([]byte(b64.encEntry.Text))
		b64.encEntryOut.SetText(encoded)
	}
}
