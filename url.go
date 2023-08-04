package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type urlTab struct {
	content *container.TabItem

	decEntry    *widget.Entry
	decEntryOut *widget.Entry
	decButton   *widget.Button

	encEntry    *widget.Entry
	encEntryOut *widget.Entry
	encButton   *widget.Button
}

func newUrlTab() *urlTab {
	u := &urlTab{}
	u.generateUI()
	u.setButtonsHandlers()
	return u
}

func (u *urlTab) generateUI() {
	u.content = container.NewTabItem(
		"URL",
		container.New(
			layout.NewGridLayout(2),
			u.generateDecodeUI(),
			u.generateEncodeUI(),
		),
	)
}

func (u *urlTab) generateDecodeUI() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("URL decoder", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true})
	outLabel := widget.NewLabelWithStyle("Decoded string:", fyne.TextAlignLeading, fyne.TextStyle{})

	u.decEntry = widget.NewMultiLineEntry()
	u.decEntryOut = widget.NewMultiLineEntry()
	u.decButton = widget.NewButton("DECODE", func() {})

	u.decEntry.SetPlaceHolder("Enter URL encoded string here")
	u.decEntryOut.Disable()

	outLabelCont := container.NewBorder(nil, outLabel, nil, nil)
	return container.NewBorder(title, nil, nil, nil, container.New(layout.NewGridLayout(1), u.decEntry, u.decButton, outLabelCont, u.decEntryOut))
}

func (u *urlTab) generateEncodeUI() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("URL encoder", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true})
	outLabel := widget.NewLabelWithStyle("Encoded string:", fyne.TextAlignLeading, fyne.TextStyle{})

	u.encEntry = widget.NewMultiLineEntry()
	u.encEntryOut = widget.NewMultiLineEntry()
	u.encButton = widget.NewButton("ENCODE", func() {})

	u.encEntry.SetPlaceHolder("Enter string here to URL encode it")
	u.encEntryOut.Disable()

	outLabelCont := container.NewBorder(nil, outLabel, nil, nil)
	return container.NewBorder(title, nil, nil, nil, container.New(layout.NewGridLayout(1), u.encEntry, u.encButton, outLabelCont, u.encEntryOut))
}

func (u *urlTab) setButtonsHandlers() {
	// Set DECODE button handler
	u.decButton.OnTapped = func() {
		if u.decEntry.Text == "" {
			return
		}
		decoded, err := url.QueryUnescape(u.decEntry.Text)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}
		u.decEntryOut.SetText(string(decoded))
	}

	// Set ENCODE button handler
	u.encButton.OnTapped = func() {
		if u.encEntry.Text == "" {
			return
		}
		encoded := url.QueryEscape(u.encEntry.Text)
		u.encEntryOut.SetText(encoded)
	}
}
