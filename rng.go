package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type rngTab struct {
	content *container.TabItem

	minEntry       *widget.Entry
	maxEntry       *widget.Entry
	resEntry       *widget.Entry
	generateButton *widget.Button
}

func newRngTab() *rngTab {
	r := &rngTab{}
	r.generateUI()
	r.setButtonsHandlers()
	return r
}

func (r *rngTab) generateUI() {
	r.minEntry = widget.NewEntry()
	r.maxEntry = widget.NewEntry()
	r.generateButton = widget.NewButton("Generate", func() {})
	r.resEntry = widget.NewEntry()

	r.minEntry.SetText("0")
	r.maxEntry.SetText("100")
	r.resEntry.Disable()

	r.content = container.NewTabItem(
		"RNG",
		container.New(
			layout.NewGridLayout(1),
			widget.NewLabelWithStyle("Random Number Generator", fyne.TextAlignCenter, fyne.TextStyle{}),
			container.New(
				layout.NewFormLayout(),
				widget.NewLabel("Min"),
				r.minEntry,
				widget.NewLabel("Max"),
				r.maxEntry,
			),
			r.generateButton,
			widget.NewSeparator(),

			r.resEntry,
		),
	)
}

func (r *rngTab) setButtonsHandlers() {
	r.generateButton.OnTapped = func() {
		if r.minEntry.Text == "" || r.maxEntry.Text == "" {
			return
		}

		min, err := strconv.Atoi(r.minEntry.Text)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}

		max, err := strconv.Atoi(r.maxEntry.Text)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}

		generated := rand.Int63n(int64(max - min + 1))
		r.resEntry.SetText(fmt.Sprint(generated + int64(min)))
	}
}
