package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createMainWindow(a fyne.App, reader *os.File, writer *os.File) fyne.Window {
	w := a.NewWindow("BankSight")

	selectedBank := "PKO"
	chooseBankRadio := widget.NewRadioGroup([]string{"PKO", "ING"}, func(bank string) {
		selectedBank = bank
	})
	chooseBankRadio.SetSelected("PKO")

	var selectedFilePath string
	selectedFileLabel := widget.NewLabel("No file selected")
	selectedFileLabel.Wrapping = fyne.TextWrapWord

	uploadButton := widget.NewButton("Upload (or drop)", func() {
		fileOpenDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			onFileOpen(reader, w, &selectedFilePath, selectedFileLabel, err)
		}, w)

		bank, err := getBank(selectedBank)

		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		downloadsFolder, err := getDownloadsFolderLocation()

		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		fileOpenDialog.SetFilter(storage.NewExtensionFileFilter(bank.getSupportedFileExtensions()))
		fileOpenDialog.SetLocation(downloadsFolder)
		fileOpenDialog.SetView(dialog.ListView)
		fileOpenDialog.Resize(fyne.NewSize(650, 400))
		fileOpenDialog.Show()
	})

	w.SetOnDropped(func(_ fyne.Position, u []fyne.URI) {
		if len(u) == 0 {
			return
		}
		setSelectedFilePath(u[0].Path(), &selectedFilePath, selectedFileLabel)
	})

	getInsightsButton := widget.NewButton("Get insights", func() {
		if selectedFilePath == "" {
			dialog.ShowInformation("No file selected", "Please select a file to upload.", w)
			return
		}

		getInsights(selectedBank, selectedFilePath, w)
	})

	terminal := widget.NewMultiLineEntry()
	terminal.SetMinRowsVisible(8)
	terminal.Wrapping = fyne.TextWrapWord
	terminal.SetPlaceHolder("Program output (stdout/stderr) will appear here...")
	terminal.Disable()

	listenToProgramOutput(reader, writer, terminal)

	toggleTerminalButton := widget.NewButtonWithIcon("", theme.ErrorIcon(), func() {
		toggleTerminal(terminal)
	})

	const padding = 10

	content := container.New(
		layout.NewCustomPaddedLayout(padding, padding, padding, padding),
		container.NewVBox(
			container.NewHBox(
				widget.NewLabel("Choose your bank:"),
				layout.NewSpacer(),
				toggleTerminalButton,
			),
			chooseBankRadio,
			widget.NewLabel("Upload your bank statement:"),
			container.NewBorder(
				nil,
				nil,
				uploadButton,
				nil,
				selectedFileLabel,
			),
			container.NewCenter(
				getInsightsButton,
			),
			terminal,
		),
	)

	w.Resize(fyne.NewSize(800, 0))
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.SetContent(content)

	return w
}

func onFileOpen(
	reader fyne.URIReadCloser,
	w fyne.Window,
	selectedFilePath *string,
	selectedFileLabel *widget.Label,
	err error,
) {
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	if reader == nil {
		return
	}

	setSelectedFilePath(reader.URI().Path(), selectedFilePath, selectedFileLabel)
	reader.Close()
}

func getDownloadsFolderLocation() (fyne.ListableURI, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	downloadsPath := filepath.Join(homeDir, "Downloads")
	downloadsListableURI, err := storage.ListerForURI(storage.NewFileURI(downloadsPath))

	if err != nil {
		return nil, err
	}

	return downloadsListableURI, nil
}

func setSelectedFilePath(path string, selectedFilePath *string, selectedFileLabel *widget.Label) {
	*selectedFilePath = path
	selectedFileLabel.SetText(filepath.Base(path))
}

func getInsights(selectedBank string, selectedFilePath string, w fyne.Window) {
	bank, err := getBank(selectedBank)

	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	textFilePath, err := bank.prepareFile(selectedFilePath)

	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	lines, err := readAllLines(textFilePath)

	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	transactions, err := bank.calculateTotalTransactions(lines)

	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	err = drawPieChart(transactions)

	if err != nil {
		dialog.ShowError(err, w)
		return
	}
}

func listenToProgramOutput(reader *os.File, writer *os.File, terminal *widget.Entry) {
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)

	go func() {
		defer reader.Close()

		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			line := scanner.Text()

			fyne.Do(func() {
				terminal.Text += line + "\n"
				terminal.CursorRow = len(terminal.Text) // Scroll to the bottom
				terminal.Refresh()
			})
		}
	}()
}

func toggleTerminal(terminal *widget.Entry) {
	if terminal.Visible() {
		terminal.Hide()
	} else {
		terminal.Show()
	}
}
