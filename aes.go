package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type aesTab struct {
	content *container.TabItem

	decEntryCipher *widget.Entry
	decEntryKey    *widget.Entry
	decEntryPlain  *widget.Entry
	decButton      *widget.Button

	encEntryPlain  *widget.Entry
	encEntryKey    *widget.Entry
	encEntryCipher *widget.Entry
	encButton      *widget.Button

	decFileCipherBytes        []byte
	decFileLabelCipherPath    *widget.Label
	decFileCipherBrowseButton *widget.Button
	decFileLabelPlainPath     *widget.Label
	decFilePlainBrowseButton  *widget.Button
	decFileEntryKey           *widget.Entry
	decFileButton             *widget.Button

	encFilePlainBytes         []byte
	encFileLabelPlainPath     *widget.Label
	encFilePlainBrowseButton  *widget.Button
	encFileLabelCipherPath    *widget.Label
	encFileCipherBrowseButton *widget.Button
	encFileEntryKey           *widget.Entry
	encFileButton             *widget.Button
}

func newAesTab() *aesTab {
	a := &aesTab{}
	a.generateUI()
	a.setButtonsHandlers()
	return a
}

func (a *aesTab) generateUI() {
	a.content = container.NewTabItem(
		"AES",
		widget.NewAccordion(
			widget.NewAccordionItem("AES decryption", a.generateDecryptionUI()),
			widget.NewAccordionItem("AES encryption", a.generateEncryptionUI()),
			widget.NewAccordionItem("AES file decryption", a.generateFileDecryptionUI()),
			widget.NewAccordionItem("AES file encryption", a.generateFileEncryptionUI()),
		),
	)
}

func (a *aesTab) generateDecryptionUI() fyne.CanvasObject {
	a.decEntryCipher = widget.NewMultiLineEntry()
	a.decEntryKey = widget.NewEntry()
	a.decEntryPlain = widget.NewMultiLineEntry()
	a.decButton = widget.NewButton("DECRYPT", func() {})
	outLabel := widget.NewLabelWithStyle("Decrypted plaintext:", fyne.TextAlignLeading, fyne.TextStyle{})

	a.decEntryCipher.SetPlaceHolder("Enter ciphertext that you want to decrypt, as a hexadecimal value")
	a.decEntryKey.SetPlaceHolder("Enter key here")
	a.decEntryPlain.Disable()

	outLabelCont := container.NewBorder(nil, outLabel, nil, nil)
	return container.New(layout.NewGridLayout(1), a.decEntryCipher, a.decEntryKey, a.decButton, outLabelCont, a.decEntryPlain)
}

func (a *aesTab) generateEncryptionUI() fyne.CanvasObject {
	a.encEntryPlain = widget.NewMultiLineEntry()
	a.encEntryKey = widget.NewEntry()
	a.encEntryCipher = widget.NewMultiLineEntry()
	a.encButton = widget.NewButton("ENCRYPT", func() {})
	outLabel := widget.NewLabelWithStyle("Encrypted ciphertext:", fyne.TextAlignLeading, fyne.TextStyle{})

	a.encEntryPlain.SetPlaceHolder("Enter plaintext that you want to encrypt")
	a.encEntryKey.SetPlaceHolder("Enter key here")
	a.encEntryCipher.Disable()

	outLabelCont := container.NewBorder(nil, outLabel, nil, nil)
	return container.New(layout.NewGridLayout(1), a.encEntryPlain, a.encEntryKey, a.encButton, outLabelCont, a.encEntryCipher)
}

func (a *aesTab) generateFileDecryptionUI() fyne.CanvasObject {
	a.decFileLabelCipherPath = widget.NewLabel("")
	a.decFileCipherBrowseButton = widget.NewButtonWithIcon("Browse", theme.FolderIcon(), func() {})
	inputFileCont := container.New(
		layout.NewGridLayout(1),
		widget.NewLabel("Input ciphertext file:"),
		container.NewBorder(nil, nil, a.decFileCipherBrowseButton, nil, a.decFileLabelCipherPath),
	)

	a.decFileEntryKey = widget.NewEntry()
	a.decFileEntryKey.SetPlaceHolder("Enter key here")

	a.decFileLabelPlainPath = widget.NewLabel("")
	a.decFilePlainBrowseButton = widget.NewButtonWithIcon("Browse", theme.FolderIcon(), func() {})
	outputFileCont := container.New(
		layout.NewGridLayout(1),
		widget.NewLabel("Output plaintext file location:"),
		container.NewBorder(nil, nil, a.decFilePlainBrowseButton, nil, a.decFileLabelPlainPath),
	)

	a.decFileButton = widget.NewButton("DECRYPT FILE", func() {})

	return container.New(layout.NewGridLayout(1), inputFileCont, a.decFileEntryKey, outputFileCont, a.decFileButton)
}

func (a *aesTab) generateFileEncryptionUI() fyne.CanvasObject {
	a.encFileLabelPlainPath = widget.NewLabel("")
	a.encFilePlainBrowseButton = widget.NewButtonWithIcon("Browse", theme.FolderIcon(), func() {})
	inputFileCont := container.New(
		layout.NewGridLayout(1),
		widget.NewLabel("Input plaintext file:"),
		container.NewBorder(nil, nil, a.encFilePlainBrowseButton, nil, a.encFileLabelPlainPath),
	)

	a.encFileEntryKey = widget.NewEntry()
	a.encFileEntryKey.SetPlaceHolder("Enter key here")

	a.encFileLabelCipherPath = widget.NewLabel("")
	a.encFileCipherBrowseButton = widget.NewButtonWithIcon("Browse", theme.FolderIcon(), func() {})
	outputFileCont := container.New(
		layout.NewGridLayout(1),
		widget.NewLabel("Output ciphertext file location:"),
		container.NewBorder(nil, nil, a.encFileCipherBrowseButton, nil, a.encFileLabelCipherPath),
	)

	a.encFileButton = widget.NewButton("ENCRYPT FILE", func() {})

	return container.New(layout.NewGridLayout(1), inputFileCont, a.encFileEntryKey, outputFileCont, a.encFileButton)
}

func (a *aesTab) setButtonsHandlers() {
	a.setDecyptTextButtonsHandlers()
	a.setEncyptTextButtonsHandlers()
	a.setDecyptFileButtonsHandlers()
	a.setEncyptFileButtonsHandlers()
}

func (a *aesTab) setDecyptTextButtonsHandlers() {
	a.decButton.OnTapped = func() {
		if a.decEntryCipher.Text == "" {
			return
		}

		cyphertext, err := hex.DecodeString(a.decEntryCipher.Text)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}

		plaintext, err := aesDecrypt(cyphertext, a.decEntryKey.Text)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}

		a.decEntryPlain.SetText(string(plaintext))
	}
}

func (a *aesTab) setEncyptTextButtonsHandlers() {
	a.encButton.OnTapped = func() {
		if a.encEntryPlain.Text == "" {
			return
		}

		ciphertext, err := aesEncrypt([]byte(a.encEntryPlain.Text), a.encEntryKey.Text)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}

		a.encEntryCipher.SetText(fmt.Sprintf("%x\n", ciphertext))
	}
}

func (a *aesTab) setDecyptFileButtonsHandlers() {
	// Set BROWSE CIPHERTEXT FILE button handler
	a.decFileCipherBrowseButton.OnTapped = func() {
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
			a.decFileCipherBytes = data
			a.decFileLabelCipherPath.SetText(reader.URI().Path())
			a.decFileLabelPlainPath.SetText(reader.URI().Path())
		}, mainWindow)

		dialog.Show()
	}

	// Set BROWSE PLAINTEXT FILE FOLDER button handler
	a.decFilePlainBrowseButton.OnTapped = func() {
		dialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			if uri == nil {
				return
			}

			if a.decFileLabelCipherPath.Text == "" {
				a.decFileLabelPlainPath.SetText(uri.Path())
				return
			}
			a.decFileLabelPlainPath.SetText(path.Join(uri.Path(), path.Base(a.decFileLabelCipherPath.Text)))
		}, mainWindow)

		dialog.Show()
	}

	// Set DECRYPT FILE button handler
	a.decFileButton.OnTapped = func() {
		plaintext, err := aesDecrypt(a.decFileCipherBytes, a.decFileEntryKey.Text)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}

		if err = os.RemoveAll(a.decFileLabelPlainPath.Text); err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}
		if err = os.WriteFile(a.decFileLabelPlainPath.Text, plaintext, 0664); err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}
		dialog.ShowInformation("File saved", fmt.Sprintf("File successfully saved to %s", a.decFileLabelPlainPath.Text), mainWindow)
	}
}

func (a *aesTab) setEncyptFileButtonsHandlers() {
	// Set BROWSE PLAINTEXT FILE button handler
	a.encFilePlainBrowseButton.OnTapped = func() {
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
			a.encFilePlainBytes = data
			a.encFileLabelPlainPath.SetText(reader.URI().Path())
			a.encFileLabelCipherPath.SetText(reader.URI().Path())
		}, mainWindow)

		dialog.Show()
	}

	// Set BROWSE CIPHERTEXT FILE FOLDER button handler
	a.encFileCipherBrowseButton.OnTapped = func() {
		dialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			if uri == nil {
				return
			}

			if a.encFileLabelPlainPath.Text == "" {
				a.encFileLabelCipherPath.SetText(uri.Path())
				return
			}
			a.encFileLabelCipherPath.SetText(path.Join(uri.Path(), path.Base(a.encFileLabelPlainPath.Text)))
		}, mainWindow)

		dialog.Show()
	}

	// Set ENCRYPT FILE button handler
	a.encFileButton.OnTapped = func() {
		ciphertext, err := aesEncrypt(a.encFilePlainBytes, a.encFileEntryKey.Text)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}

		if err = os.RemoveAll(a.encFileLabelCipherPath.Text); err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}
		if err = os.WriteFile(a.encFileLabelCipherPath.Text, ciphertext, 0664); err != nil {
			dialog.ShowError(err, mainWindow)
			return
		}
		dialog.ShowInformation("File saved", fmt.Sprintf("File successfully saved to %s", a.encFileLabelCipherPath.Text), mainWindow)
	}
}

func aesDecrypt(ciphertextBytes []byte, keyString string) ([]byte, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertextBytes) < aes.BlockSize {
		return nil, errors.New("ciphertext is too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]
	plaintextBytes := make([]byte, len(ciphertextBytes))

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintextBytes, ciphertextBytes)

	return plaintextBytes, nil
}

func aesEncrypt(plaintextBytes []byte, keyString string) ([]byte, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	return ciphertext, nil
}
