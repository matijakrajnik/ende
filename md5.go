package main

import (
	"crypto/md5"
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type md5Tab struct {
	content *container.TabItem

	inputStringEntry  *widget.Entry
	outputStringEntry *widget.Entry
	hashStringButton  *widget.Button

	inputFileBytes        []byte
	inputFileLabel        *widget.Label
	inputFileBrowseButton *widget.Button
	outputFileEntry       *widget.Entry
	hashFileButton        *widget.Button
}

func newMd5Tab() *md5Tab {
	m := &md5Tab{}
	m.generateUI()
	m.setButtonsHandlers()
	return m
}

func (m *md5Tab) generateUI() {
	m.content = container.NewTabItem(
		"MD5",
		container.New(
			layout.NewVBoxLayout(),
			widget.NewLabelWithStyle("MD5 Hash", fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewAccordion(
				widget.NewAccordionItem("String hash", m.generateStringHashUI()),
				widget.NewAccordionItem("File hash", m.generateFileHashUI()),
			),
		),
	)
}

func (m *md5Tab) generateStringHashUI() fyne.CanvasObject {
	m.inputStringEntry = widget.NewEntry()
	m.outputStringEntry = widget.NewEntry()
	m.hashStringButton = widget.NewButton("CALCULATE HASH", func() {})
	hashLabel := widget.NewLabelWithStyle("Hash:", fyne.TextAlignLeading, fyne.TextStyle{})

	m.inputStringEntry.SetPlaceHolder("Enter your string here")
	m.outputStringEntry.Disable()

	hashLabelCont := container.NewBorder(nil, hashLabel, nil, nil)
	return container.New(layout.NewGridLayout(1), m.inputStringEntry, m.hashStringButton, hashLabelCont, m.outputStringEntry)
}

func (m *md5Tab) generateFileHashUI() fyne.CanvasObject {
	m.inputFileLabel = widget.NewLabel("")
	m.inputFileBrowseButton = widget.NewButtonWithIcon("Browse", theme.FolderIcon(), func() {})
	inputFileCont := container.New(
		layout.NewGridLayout(1),
		widget.NewLabel("Input file:"),
		container.NewBorder(nil, nil, m.inputFileBrowseButton, nil, m.inputFileLabel),
	)

	m.outputFileEntry = widget.NewEntry()
	m.hashFileButton = widget.NewButton("CALCULATE HASH", func() {})
	hashLabel := widget.NewLabelWithStyle("Hash:", fyne.TextAlignLeading, fyne.TextStyle{})

	m.outputFileEntry.Disable()

	hashLabelCont := container.NewBorder(nil, hashLabel, nil, nil)
	return container.New(layout.NewGridLayout(1), inputFileCont, m.hashFileButton, hashLabelCont, m.outputFileEntry)
}

func (m *md5Tab) setButtonsHandlers() {
	m.setStringHashButtonsHandlers()
	m.setFileHashButtonsHandlers()
}

func (m *md5Tab) setStringHashButtonsHandlers() {
	m.hashStringButton.OnTapped = func() {
		if m.inputStringEntry.Text == "" {
			return
		}
		m.outputStringEntry.SetText(fmt.Sprintf("%x", md5.Sum([]byte(m.inputStringEntry.Text))))
	}
}

func (m *md5Tab) setFileHashButtonsHandlers() {
	m.inputFileBrowseButton.OnTapped = func() {
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			if reader == nil {
				return
			}
			data, err := io.ReadAll(reader)
			defer reader.Close()
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			m.inputFileBytes = data
			m.inputFileLabel.SetText(reader.URI().Path())
		}, mainWindow)

		dialog.Show()
	}

	m.hashFileButton.OnTapped = func() {
		if m.inputFileLabel.Text == "" {
			return
		}
		m.outputFileEntry.SetText(fmt.Sprintf("%x", md5.Sum(m.inputFileBytes)))
	}
}
